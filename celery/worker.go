package celery

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/cloud-mill/go-celery/logger"
	"github.com/cloud-mill/go-celery/models"
	"go.uber.org/zap"
)

type WorkerPool struct {
	broker     Broker
	backend    Backend
	numWorkers int
	tasks      map[string]Task
	rateLimit  time.Duration
	lock       sync.RWMutex
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

func NewWorkerPool(
	broker Broker,
	backend Backend,
	numWorkers int,
) *WorkerPool {
	return &WorkerPool{
		broker:     broker,
		backend:    backend,
		numWorkers: numWorkers,
		tasks:      make(map[string]Task),
		rateLimit:  100 * time.Millisecond,
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	ctx, wp.cancel = context.WithCancel(ctx)
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(ctx, i)
	}
}

func (wp *WorkerPool) Stop() {
	if wp.cancel != nil {
		wp.cancel()
	}

	wp.wg.Wait()
}

func (wp *WorkerPool) RegisterTask(
	name string,
	task Task,
) {
	wp.lock.Lock()
	defer wp.lock.Unlock()

	wp.tasks[name] = task
}

func (wp *WorkerPool) getTask(name string) Task {
	wp.lock.RLock()
	defer wp.lock.RUnlock()

	task, exists := wp.tasks[name]
	if !exists {
		return nil
	}
	return task
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()
	ticker := time.NewTicker(wp.rateLimit)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			taskMessage, err := wp.broker.GetTaskMessage(ctx)
			if err != nil {
				logger.Logger.Error(
					"broker failed to get task message",
					zap.Int("worker id", id),
					zap.Error(err),
				)
				continue
			}

			if taskMessage == nil {
				continue
			}

			receivedEvent := models.TaskReceivedEvent{
				BaseEvent: models.BaseEvent{
					Event:     models.EventTaskReceived,
					UUID:      taskMessage.Headers.ID,
					Timestamp: time.Now(),
					Type:      taskMessage.Headers.Task,
					Hostname:  "worker" + strconv.Itoa(id),
				},
				Args:   taskMessage.Body.Args,
				Kwargs: taskMessage.Body.Kwargs,
			}
			_ = wp.broker.SendTaskReceivedEvent(ctx, receivedEvent)

			startedEvent := models.TaskStartedEvent{
				BaseEvent: models.BaseEvent{
					Event:     models.EventTaskStarted,
					UUID:      taskMessage.Headers.ID,
					Timestamp: time.Now(),
					Type:      taskMessage.Headers.Task,
					Hostname:  "worker" + strconv.Itoa(id),
				},
			}
			_ = wp.broker.SendTaskStartedEvent(ctx, startedEvent)

			result, err := wp.executeTask(*taskMessage)
			if err != nil {
				logger.Logger.Error(
					"worker failed to run task",
					zap.String("task message id", taskMessage.Headers.ID),
					zap.Error(err),
				)

				failedEvent := models.TaskFailedEvent{
					BaseEvent: models.BaseEvent{
						Event:     models.EventTaskFailed,
						UUID:      taskMessage.Headers.ID,
						Timestamp: time.Now(),
						Type:      taskMessage.Headers.Task,
						Hostname:  "worker" + strconv.Itoa(id),
					},
					Exception: err.Error(),
					Traceback: "",
				}
				_ = wp.broker.SendTaskFailedEvent(ctx, failedEvent)
				continue
			}

			if err = wp.backend.SetResult(ctx, taskMessage.Headers.ID, result); err != nil {
				logger.Logger.Error(
					"failed to push task result",
					zap.String("taskMessageId", taskMessage.Headers.ID),
					zap.Error(err),
				)
				continue
			}

			succeededEvent := models.TaskSucceededEvent{
				BaseEvent: models.BaseEvent{
					Event:     models.EventTaskSucceeded,
					UUID:      taskMessage.Headers.ID,
					Timestamp: time.Now(),
					Type:      taskMessage.Headers.Task,
					Hostname:  "worker" + strconv.Itoa(id),
				},
				Result: result,
			}
			_ = wp.broker.SendTaskSucceededEvent(ctx, succeededEvent)
		}
	}
}

func (wp *WorkerPool) executeTask(message models.TaskMessage) (interface{}, error) {
	task := wp.getTask(message.Headers.Task)
	if task == nil {
		return nil, fmt.Errorf("task %s is not registered", message.Headers.Task)
	}

	if err := task.ParseKwargs(message.Body.Kwargs); err != nil {
		return nil, err
	}
	val, err := task.RunTask()
	if err != nil {
		return nil, err
	}

	return val, nil
}

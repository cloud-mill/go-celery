package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloud-mill/go-celery/logger"
	"time"

	"github.com/cloud-mill/go-celery/models"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type CeleryRedisBroker struct {
	RedisClient    *redis.Client
	QueueName      string
	EventQueueName string
}

func NewCeleryRedisBroker(redisClient *redis.Client) *CeleryRedisBroker {
	return &CeleryRedisBroker{
		RedisClient:    redisClient,
		QueueName:      "go-celery",
		EventQueueName: "go-celery-events",
	}
}

func (broker *CeleryRedisBroker) SendTaskMessage(
	ctx context.Context,
	taskMessage models.TaskMessage,
) error {
	jsonBytes, err := json.Marshal(taskMessage)
	if err != nil {
		logger.Logger.Error("failed to marshal task message", zap.Error(err))
		return err
	}

	_, err = broker.RedisClient.LPush(ctx, broker.QueueName, jsonBytes).Result()
	if err != nil {
		logger.Logger.Error("failed to push task message to Redis", zap.Error(err))
		return err
	}

	return nil
}

func (broker *CeleryRedisBroker) GetTaskMessage(
	ctx context.Context,
) (*models.TaskMessage, error) {
	res, err := broker.RedisClient.BRPop(ctx, 1*time.Second, broker.QueueName).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		logger.Logger.Error("failed to pop message from Redis", zap.Error(err))
		return nil, err
	}

	if len(res) < 2 {
		err := fmt.Errorf("received an invalid message format from Redis")
		logger.Logger.Error("invalid message format", zap.Error(err))
		return nil, err
	}

	var message models.TaskMessage
	if err := json.Unmarshal([]byte(res[1]), &message); err != nil {
		logger.Logger.Error("failed to unmarshal message", zap.Error(err))
		return nil, err
	}

	return &message, nil
}

func (broker *CeleryRedisBroker) SendTaskReceivedEvent(
	ctx context.Context,
	event models.TaskReceivedEvent,
) error {
	return broker.sendEventMessage(ctx, event)
}

func (broker *CeleryRedisBroker) SendTaskStartedEvent(
	ctx context.Context,
	event models.TaskStartedEvent,
) error {
	return broker.sendEventMessage(ctx, event)
}

func (broker *CeleryRedisBroker) SendTaskSucceededEvent(
	ctx context.Context,
	event models.TaskSucceededEvent,
) error {
	return broker.sendEventMessage(ctx, event)
}

func (broker *CeleryRedisBroker) SendTaskFailedEvent(
	ctx context.Context,
	event models.TaskFailedEvent,
) error {
	return broker.sendEventMessage(ctx, event)
}

func (broker *CeleryRedisBroker) sendEventMessage(ctx context.Context, event interface{}) error {
	jsonBytes, err := json.Marshal(event)
	if err != nil {
		logger.Logger.Error("failed to marshal event message", zap.Error(err))
		return err
	}

	_, err = broker.RedisClient.LPush(ctx, broker.EventQueueName, jsonBytes).Result()
	if err != nil {
		logger.Logger.Error("failed to push event message to Redis", zap.Error(err))
		return err
	}

	return nil
}

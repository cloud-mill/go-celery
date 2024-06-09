package celery

import (
	"context"

	"github.com/cloud-mill/go-celery/models"
)

type Broker interface {
	SendTaskMessage(
		ctx context.Context,
		taskMessage models.TaskMessage,
	) error
	GetTaskMessage(ctx context.Context) (*models.TaskMessage, error)

	SendTaskReceivedEvent(ctx context.Context, event models.TaskReceivedEvent) error
	SendTaskStartedEvent(ctx context.Context, event models.TaskStartedEvent) error
	SendTaskSucceededEvent(ctx context.Context, event models.TaskSucceededEvent) error
	SendTaskFailedEvent(ctx context.Context, event models.TaskFailedEvent) error
}

type Backend interface {
	GetResult(ctx context.Context, taskID string) (interface{}, error)
	SetResult(ctx context.Context, taskID string, result interface{}) error
}

type Task interface {
	ParseKwargs(map[string]interface{}) error
	RunTask() (interface{}, error)
	HandleEvent(ctx context.Context, event models.BaseEvent) error
}

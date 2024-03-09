package celery

import (
	"context"

	"github.com/cloud-mill/go-celery/models"
)

type Broker interface {
	SendCeleryMessage(
		ctx context.Context,
		message models.CeleryMessage,
	) error
	GetTaskMessage(ctx context.Context) (*models.TaskMessage, error)
}

type Backend interface {
	GetResult(ctx context.Context, taskId string) (*models.ResultMessage, error)
	SetResult(ctx context.Context, taskId string, result models.ResultMessage) error
}

type Task interface {
	ParseKwargs(map[string]interface{}) error
	RunTask() (interface{}, error)
}

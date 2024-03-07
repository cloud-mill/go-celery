package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cloud-mill/go-celery/models"
	"github.com/go-redis/redis/v8"
)

type CeleryBackend struct {
	RedisClient *redis.Client
}

func NewRedisBackend(redisClient *redis.Client) *CeleryBackend {
	return &CeleryBackend{
		RedisClient: redisClient,
	}
}

func (celeryBackend *CeleryBackend) GetResult(
	ctx context.Context,
	taskId string,
) (*models.ResultMessage, error) {
	val, err := celeryBackend.RedisClient.Get(ctx, taskId).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("result not available for task Id: %s", taskId)
	} else if err != nil {
		return nil, err
	}

	var resultMessage models.ResultMessage
	if err := json.Unmarshal(val, &resultMessage); err != nil {
		return nil, fmt.Errorf("error unmarshalling result for task Id %s: %w", taskId, err)
	}

	return &resultMessage, nil
}

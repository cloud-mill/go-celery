package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloud-mill/go-celery/logger"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type CeleryRedisBackend struct {
	RedisClient *redis.Client
}

func NewCeleryRedisBackend(redisClient *redis.Client) *CeleryRedisBackend {
	return &CeleryRedisBackend{
		RedisClient: redisClient,
	}
}

func (backend *CeleryRedisBackend) GetResult(
	ctx context.Context,
	taskId string,
) (interface{}, error) {
	val, err := backend.RedisClient.Get(ctx, taskId).Bytes()
	if errors.Is(err, redis.Nil) {
		err := fmt.Errorf("result not available for task Id: %s", taskId)
		logger.Logger.Error("result not available", zap.String("taskId", taskId), zap.Error(err))
		return nil, err
	} else if err != nil {
		logger.Logger.Error("failed to get result", zap.String("taskId", taskId), zap.Error(err))
		return nil, err
	}

	return val, nil
}

func (backend *CeleryRedisBackend) SetResult(
	ctx context.Context,
	taskId string,
	result interface{},
) error {
	err := backend.RedisClient.Set(ctx, taskId, result, 0).Err()
	if err != nil {
		logger.Logger.Error("failed to set result", zap.String("taskId", taskId), zap.Error(err))
	}
	return err
}

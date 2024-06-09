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
	taskID string,
) (interface{}, error) {
	val, err := backend.RedisClient.Get(ctx, taskID).Bytes()
	if errors.Is(err, redis.Nil) {
		err := fmt.Errorf("result not available for task Id: %s", taskID)
		logger.Logger.Error("result not available", zap.String("taskID", taskID), zap.Error(err))
		return nil, err
	} else if err != nil {
		logger.Logger.Error("failed to get result", zap.String("taskID", taskID), zap.Error(err))
		return nil, err
	}

	return val, nil
}

func (backend *CeleryRedisBackend) SetResult(
	ctx context.Context,
	taskID string,
	result interface{},
) error {
	err := backend.RedisClient.Set(ctx, taskID, result, 0).Err()
	if err != nil {
		logger.Logger.Error("failed to set result", zap.String("taskID", taskID), zap.Error(err))
	}
	return err
}

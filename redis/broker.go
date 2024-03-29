package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cloud-mill/go-celery/models"
	"github.com/go-redis/redis/v8"
)

type CeleryRedisBroker struct {
	RedisClient *redis.Client
	QueueName   string
}

func NewRedisCeleryRedisBroker(redisClient *redis.Client) *CeleryRedisBroker {
	return &CeleryRedisBroker{
		RedisClient: redisClient,
		QueueName:   "go-celery",
	}
}

func (celeryRedisBroker *CeleryRedisBroker) SendCeleryMessage(
	ctx context.Context,
	message models.CeleryMessage,
) error {
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = celeryRedisBroker.RedisClient.LPush(ctx, celeryRedisBroker.QueueName, jsonBytes).
		Result()
	if err != nil {
		return err
	}

	return nil
}

func (celeryRedisBroker *CeleryRedisBroker) getCeleryMessage(
	ctx context.Context,
) (*models.CeleryMessage, error) {
	// BRPOP command to pop the last(right) message from the list (queue), with a timeout of 1 second
	res, err := celeryRedisBroker.RedisClient.BRPop(ctx, 1*time.Second, celeryRedisBroker.QueueName).
		Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	if len(res) < 2 {
		// The first element in res is the queue name, and the second is the message.
		// If we have less than 2 elements, something went wrong.
		return nil, fmt.Errorf("received an invalid message format from Redis")
	}

	var message models.CeleryMessage
	if err := json.Unmarshal([]byte(res[1]), &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (celeryRedisBroker *CeleryRedisBroker) GetTaskMessage(
	ctx context.Context,
) (*models.TaskMessage, error) {
	celeryMessage, err := celeryRedisBroker.getCeleryMessage(ctx)
	if err != nil || celeryMessage == nil {
		return nil, err
	}

	return celeryMessage.ExtractTaskMessage()
}

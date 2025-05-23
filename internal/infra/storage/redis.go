package storage

import (
	"context"
	"errors"
	"kafka-connect/internal/application/repo"
	"kafka-connect/internal/infra/config"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func (r Redis) Save(ctx context.Context, key string, val []byte) error {
	return r.client.Set(ctx, key, val, 0).Err()
}

func (r Redis) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	notExists := err != nil && errors.Is(err, redis.Nil)
	if notExists {
		return data, repo.NotExistsKey
	}
	return data, err
}

func NewRedis(
	conf config.Config,
) *Redis {
	opt, _ := redis.ParseURL(conf.RedisOpt)
	client := redis.NewClient(opt)
	return &Redis{client}
}

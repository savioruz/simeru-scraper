package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.New("failed to connect to redis")
	}

	return &RedisCache{
		client: client,
	}, nil
}

func (r *RedisCache) Get(key string, value interface{}) error {
	data, err := r.client.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return errors.New("cache miss for key: " + key)
	} else if err != nil {
		return errors.New("failed to get data from redis for key: " + key)
	}

	if err := json.Unmarshal([]byte(data), value); err != nil {
		return errors.New("failed to unmarshal data")
	}

	return nil
}

func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.New("failed to marshal data for key: " + key)
	}

	if _, err := r.client.Set(context.Background(), key, data, expiration).Result(); err != nil {
		return errors.New("failed to set data to redis for key: " + key)
	}

	return nil
}

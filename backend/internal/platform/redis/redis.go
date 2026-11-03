package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(cfg *config.RedisConfig) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &Redis{client: client}
}

func (r *Redis) Client() *redis.Client {
	return r.client
}

func (r *Redis) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *Redis) healthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.Ping(ctx)
}

func (r *Redis) CheckConnection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis no está disponible: %w", err)
	}
	return nil
}

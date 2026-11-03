package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	recuperacionTokenTTL = 15 * time.Minute
	recuperacionKeyPrefix = "auth:recuperacion:"
)

type TokenStorage struct {
	client *redis.Client
}

func NewTokenStorage(client *redis.Client) *TokenStorage {
	return &TokenStorage{client: client}
}

// GuardarToken guarda un token de recuperación con TTL
func (ts *TokenStorage) GuardarToken(ctx context.Context, cedula string, token string) error {
	key := recuperacionKeyPrefix + cedula

	// Lua script para guardar token de forma atómica
	script := `
		local key = KEYS[1]
		local token = ARGV[1]
		local ttl = tonumber(ARGV[2])
		
		redis.call('SET', key, token)
		redis.call('EXPIRE', key, ttl)
		
		return 1
	`

	_, err := ts.client.Eval(ctx, script, []string{key}, token, int(recuperacionTokenTTL.Seconds())).Result()
	if err != nil {
		return fmt.Errorf("error guardando token: %w", err)
	}

	return nil
}

// ValidarToken valida un token de recuperación
func (ts *TokenStorage) ValidarToken(ctx context.Context, cedula string, token string) (bool, error) {
	key := recuperacionKeyPrefix + cedula

	script := `
		local key = KEYS[1]
		local expectedToken = ARGV[1]
		
		local storedToken = redis.call('GET', key)
		
		if storedToken == false then
			return 0
		end
		
		if storedToken == expectedToken then
			return 1
		else
			return 0
		end
	`

	result, err := ts.client.Eval(ctx, script, []string{key}, token).Int()
	if err != nil {
		return false, fmt.Errorf("error validando token: %w", err)
	}

	return result == 1, nil
}

// EliminarToken elimina un token de recuperación (después de usarlo)
func (ts *TokenStorage) EliminarToken(ctx context.Context, cedula string) error {
	key := recuperacionKeyPrefix + cedula

	_, err := ts.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error eliminando token: %w", err)
	}

	return nil
}

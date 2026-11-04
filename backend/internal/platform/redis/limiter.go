package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	intentosTTLSeconds = 900 // 15 minutos (tiempo de vida del contador)

	// Prefijos para registro
	regIntentosKeyPrefix = "auth:reg:intentos:"
	regBloqueoKeyPrefix  = "auth:reg:bloqueo:"

	// Prefijos para login
	loginIntentosKeyPrefix = "auth:login:intentos:"
	loginBloqueoKeyPrefix  = "auth:login:bloqueo:"
)

type IntentosLimiter struct {
	client            *redis.Client
	intentosKeyPrefix string
	bloqueoKeyPrefix  string
}

func NewRegistroIntentosLimiter(client *redis.Client) *IntentosLimiter {
	return &IntentosLimiter{
		client:            client,
		intentosKeyPrefix: regIntentosKeyPrefix,
		bloqueoKeyPrefix:  regBloqueoKeyPrefix,
	}
}

func NewLoginIntentosLimiter(client *redis.Client) *IntentosLimiter {
	return &IntentosLimiter{
		client:            client,
		intentosKeyPrefix: loginIntentosKeyPrefix,
		bloqueoKeyPrefix:  loginBloqueoKeyPrefix,
	}
}

// IncrementarIntento incrementa el contador de intentos de forma atómica
func (l *IntentosLimiter) IncrementarIntento(ctx context.Context, cedula string) (int, error) {
	key := l.intentosKeyPrefix + cedula
	result, err := l.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("error incrementando intentos: %w", err)
	}
	// Si es el primer intento, establecer TTL
	if result == 1 {
		if err := l.client.Expire(ctx, key, time.Duration(intentosTTLSeconds)*time.Second).Err(); err != nil {
			return 0, fmt.Errorf("error estableciendo TTL de intentos: %w", err)
		}
	}
	return int(result), nil
}

// VerificarBloqueo verifica si la cédula está bloqueada
func (l *IntentosLimiter) VerificarBloqueo(ctx context.Context, cedula string) (bool, int, error) {
	key := l.bloqueoKeyPrefix + cedula
	exists, err := l.client.Exists(ctx, key).Result()
	if err != nil {
		return false, 0, fmt.Errorf("error verificando bloqueo: %w", err)
	}
	if exists == 0 {
		return false, 0, nil
	}
	ttl, err := l.client.TTL(ctx, key).Result()
	if err != nil {
		return false, 0, fmt.Errorf("error obteniendo TTL: %w", err)
	}
	// Convertir a segundos
	return true, int(ttl.Seconds()), nil
}

// ReiniciarIntentos reinicia el contador de intentos
func (l *IntentosLimiter) ReiniciarIntentos(ctx context.Context, cedula string) error {
	key := l.intentosKeyPrefix + cedula
	return l.client.Del(ctx, key).Err()
}

// EstablecerBloqueo establece un bloqueo para la cédula
func (l *IntentosLimiter) EstablecerBloqueo(ctx context.Context, cedula string, ttl int) error {
	key := l.bloqueoKeyPrefix + cedula
	intentosKey := l.intentosKeyPrefix + cedula
	pipe := l.client.Pipeline()
	pipe.Set(ctx, key, "1", time.Duration(ttl)*time.Second)
	pipe.Del(ctx, intentosKey)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error estableciendo bloqueo: %w", err)
	}
	return nil
}

// ObtenerIntentos obtiene el número actual de intentos (para debugging/testing)
func (l *IntentosLimiter) ObtenerIntentos(ctx context.Context, cedula string) (int, error) {
	key := l.intentosKeyPrefix + cedula
	val, err := l.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("error obteniendo intentos: %w", err)
	}
	return val, nil
}

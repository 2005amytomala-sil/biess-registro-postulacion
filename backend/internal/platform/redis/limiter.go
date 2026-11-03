package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	intentosMaximos        = 3
	bloqueoTTLSeconds      = 600 // 10 minutos
	intentosTTLSeconds     = 900 // 15 minutos (tiempo de vida del contador)
	intentosKeyPrefix      = "auth:intentos:"
	bloqueoKeyPrefix       = "auth:bloqueo:"
)

type IntentosLimiter struct {
	client *redis.Client
}

func NewIntentosLimiter(client *redis.Client) *IntentosLimiter {
	return &IntentosLimiter{client: client}
}

// IncrementarIntento incrementa el contador de intentos de forma atómica
func (l *IntentosLimiter) IncrementarIntento(ctx context.Context, cedula string) (int, error) {
	key := intentosKeyPrefix + cedula
	result, err := l.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("error incrementando intentos: %w", err)
	}
	// Si es el primer intento, establecer TTL
	if result == 1 {
		l.client.Expire(ctx, key, time.Duration(intentosTTLSeconds)*time.Second)
	}
	return int(result), nil
}

// VerificarBloqueo verifica si la cédula está bloqueada
func (l *IntentosLimiter) VerificarBloqueo(ctx context.Context, cedula string) (bool, int, error) {
	key := bloqueoKeyPrefix + cedula
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
	key := intentosKeyPrefix + cedula
	return l.client.Del(ctx, key).Err()
}

// EstablecerBloqueo establece un bloqueo para la cédula
func (l *IntentosLimiter) EstablecerBloqueo(ctx context.Context, cedula string, ttl int) error {
	key := bloqueoKeyPrefix + cedula
	intentosKey := intentosKeyPrefix + cedula
	pipe := l.client.Pipeline()
	pipe.Set(ctx, key, "1", 0)
	pipe.Expire(ctx, key, time.Duration(ttl)*time.Second)
	pipe.Del(ctx, intentosKey)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error estableciendo bloqueo: %w", err)
	}
	return nil
}

// VerificarYBloquearSiNecesario verifica intentos y bloquea si es necesario
func (l *IntentosLimiter) VerificarYBloquearSiNecesario(ctx context.Context, cedula string) error {
	isBlocked, retryAfter, err := l.VerificarBloqueo(ctx, cedula)
	if err != nil {
		return err
	}
	if isBlocked {
		return &domain.AppError{
			Code:    domain.ErrBloqueado.Code,
			Message: fmt.Sprintf("%s (tiempo restante: %d segundos)", domain.ErrBloqueado.Message, retryAfter),
			Err:     domain.ErrBloqueado,
		}
	}
	intentos, err := l.IncrementarIntento(ctx, cedula)
	if err != nil {
		return err
	}
	if intentos >= intentosMaximos {
		if err := l.EstablecerBloqueo(ctx, cedula, bloqueoTTLSeconds); err != nil {
			return err
		}
		return domain.ErrBloqueado
	}
	return nil
}

// ObtenerIntentos obtiene el número actual de intentos (para debugging/testing)
func (l *IntentosLimiter) ObtenerIntentos(ctx context.Context, cedula string) (int, error) {
	key := intentosKeyPrefix + cedula
	val, err := l.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("error obteniendo intentos: %w", err)
	}
	return val, nil
}

package redis

import (
	"context"
	"fmt"

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
// Retorna el nuevo valor del contador
func (l *IntentosLimiter) IncrementarIntento(ctx context.Context, cedula string) (int, error) {
	key := intentosKeyPrefix + cedula

	// Lua script para incrementar y verificar si debe bloquear
	script := `
		local key = KEYS[1]
		local max = tonumber(ARGV[1])
		local ttl = tonumber(ARGV[2])
		
		local current = redis.call('INCR', key)
		
		-- Si es el primer intento, establecer TTL
		if current == 1 then
			redis.call('EXPIRE', key, ttl)
		end
		
		return current
	`

	result, err := l.client.Eval(ctx, script, []string{key}, intentosMaximos, intentosTTLSeconds).Int()
	if err != nil {
		return 0, fmt.Errorf("error incrementando intentos: %w", err)
	}

	return result, nil
}

// VerificarBloqueo verifica si la cédula está bloqueada
// Retorna (isBlocked, retryAfter, error)
func (l *IntentosLimiter) VerificarBloqueo(ctx context.Context, cedula string) (bool, int, error) {
	key := bloqueoKeyPrefix + cedula

	// Lua script para verificar bloqueo y obtener TTL restante
	script := `
		local key = KEYS[1]
		local exists = redis.call('EXISTS', key)
		
		if exists == 1 then
			local ttl = redis.call('TTL', key)
			return {1, ttl}
		else
			return {0, 0}
		end
	`

	result, err := l.client.Eval(ctx, script, []string{key}).Slice()
	if err != nil {
		return false, 0, fmt.Errorf("error verificando bloqueo: %w", err)
	}

	isBlocked := result[0].(int64) == 1
	retryAfter := int(result[1].(int64))

	return isBlocked, retryAfter, nil
}

// ReiniciarIntentos reinicia el contador de intentos (cuando el usuario acierta)
func (l *IntentosLimiter) ReiniciarIntentos(ctx context.Context, cedula string) error {
	key := intentosKeyPrefix + cedula

	// Lua script para reiniciar contador de forma atómica
	script := `
		local key = KEYS[1]
		redis.call('DEL', key)
		return 1
	`

	_, err := l.client.Eval(ctx, script, []string{key}).Result()
	if err != nil {
		return fmt.Errorf("error reiniciando intentos: %w", err)
	}

	return nil
}

// EstablecerBloqueo establece un bloqueo para la cédula
func (l *IntentosLimiter) EstablecerBloqueo(ctx context.Context, cedula string, ttl int) error {
	key := bloqueoKeyPrefix + cedula

	// Lua script para establecer bloqueo de forma atómica
	script := `
		local key = KEYS[1]
		local ttl = tonumber(ARGV[1])
		
		redis.call('SET', key, '1')
		redis.call('EXPIRE', key, ttl)
		
		-- También limpiar el contador de intentos
		local intentosKey = KEYS[2]
		redis.call('DEL', intentosKey)
		
		return 1
	`

	intentosKey := intentosKeyPrefix + cedula
	_, err := l.client.Eval(ctx, script, []string{key, intentosKey}, ttl).Result()
	if err != nil {
		return fmt.Errorf("error estableciendo bloqueo: %w", err)
	}

	return nil
}

// VerificarYBloquearSiNecesario verifica intentos y bloquea si es necesario
// Método de conveniencia que combina verificación de bloqueo e incremento
func (l *IntentosLimiter) VerificarYBloquearSiNecesario(ctx context.Context, cedula string) error {
	// Primero verificar si ya está bloqueado
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

	// Incrementar intentos
	intentos, err := l.IncrementarIntento(ctx, cedula)
	if err != nil {
		return err
	}

	// Si alcanzó el máximo, establecer bloqueo
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

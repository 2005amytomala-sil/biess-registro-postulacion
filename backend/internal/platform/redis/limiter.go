package redis

import (
	"context"
	"fmt"
	"sync"
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

// In-memory implementation for testing without Redis
type InMemoryIntentosLimiter struct {
	mu              sync.RWMutex
	intentos        map[string]int
	bloqueos        map[string]time.Time
	bloqueoDurations map[string]time.Duration
}

func NewInMemoryIntentosLimiter() *InMemoryIntentosLimiter {
	return &InMemoryIntentosLimiter{
		intentos:        make(map[string]int),
		bloqueos:        make(map[string]time.Time),
		bloqueoDurations: make(map[string]time.Duration),
	}
}

func (l *InMemoryIntentosLimiter) IncrementarIntento(ctx context.Context, cedula string) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.intentos[cedula]++
	return l.intentos[cedula], nil
}

func (l *InMemoryIntentosLimiter) VerificarBloqueo(ctx context.Context, cedula string) (bool, int, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	bloqueadoHasta, existe := l.bloqueos[cedula]
	if !existe {
		return false, 0, nil
	}

	if time.Now().Before(bloqueadoHasta) {
		retryAfter := int(time.Until(bloqueadoHasta).Seconds())
		return true, retryAfter, nil
	}

	// Bloqueo expirado, limpiar
	l.mu.Lock()
	delete(l.bloqueos, cedula)
	delete(l.bloqueoDurations, cedula)
	l.mu.Unlock()

	return false, 0, nil
}

func (l *InMemoryIntentosLimiter) ReiniciarIntentos(ctx context.Context, cedula string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.intentos, cedula)
	return nil
}

func (l *InMemoryIntentosLimiter) EstablecerBloqueo(ctx context.Context, cedula string, ttl int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	bloqueadoHasta := time.Now().Add(time.Duration(ttl) * time.Second)
	l.bloqueos[cedula] = bloqueadoHasta
	l.bloqueoDurations[cedula] = time.Duration(ttl) * time.Second
	delete(l.intentos, cedula)

	return nil
}

type IntentosLimiter struct {
	client         *redis.Client
	inMemoryLimiter *InMemoryIntentosLimiter
	useInMemory    bool
}

func NewIntentosLimiter(client *redis.Client) *IntentosLimiter {
	return &IntentosLimiter{
		client:         client,
		inMemoryLimiter: NewInMemoryIntentosLimiter(),
		useInMemory:    true, // Usar in-memory por defecto para pruebas
	}
}

func (l *IntentosLimiter) IncrementarIntento(ctx context.Context, cedula string) (int, error) {
	if l.useInMemory {
		return l.inMemoryLimiter.IncrementarIntento(ctx, cedula)
	}
	return l.redisIncrementarIntento(ctx, cedula)
}

func (l *IntentosLimiter) VerificarBloqueo(ctx context.Context, cedula string) (bool, int, error) {
	if l.useInMemory {
		return l.inMemoryLimiter.VerificarBloqueo(ctx, cedula)
	}
	return l.redisVerificarBloqueo(ctx, cedula)
}

func (l *IntentosLimiter) ReiniciarIntentos(ctx context.Context, cedula string) error {
	if l.useInMemory {
		return l.inMemoryLimiter.ReiniciarIntentos(ctx, cedula)
	}
	return l.redisReiniciarIntentos(ctx, cedula)
}

func (l *IntentosLimiter) EstablecerBloqueo(ctx context.Context, cedula string, ttl int) error {
	if l.useInMemory {
		return l.inMemoryLimiter.EstablecerBloqueo(ctx, cedula, ttl)
	}
	return l.redisEstablecerBloqueo(ctx, cedula, ttl)
}

// Redis implementation methods
func (l *IntentosLimiter) redisIncrementarIntento(ctx context.Context, cedula string) (int, error) {
	key := intentosKeyPrefix + cedula

	script := `
		local key = KEYS[1]
		local max = tonumber(ARGV[1])
		local ttl = tonumber(ARGV[2])
		
		local current = redis.call('INCR', key)
		
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

func (l *IntentosLimiter) redisVerificarBloqueo(ctx context.Context, cedula string) (bool, int, error) {
	key := bloqueoKeyPrefix + cedula

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

func (l *IntentosLimiter) redisReiniciarIntentos(ctx context.Context, cedula string) error {
	key := intentosKeyPrefix + cedula

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

func (l *IntentosLimiter) redisEstablecerBloqueo(ctx context.Context, cedula string, ttl int) error {
	key := bloqueoKeyPrefix + cedula

	script := `
		local key = KEYS[1]
		local ttl = tonumber(ARGV[1])
		
		redis.call('SET', key, '1')
		redis.call('EXPIRE', key, ttl)
		
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
	if l.useInMemory {
		return l.inMemoryLimiter.ObtenerIntentos(ctx, cedula)
	}

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

// ObtenerIntentos obtiene el número actual de intentos (para debugging/testing)
func (l *InMemoryIntentosLimiter) ObtenerIntentos(ctx context.Context, cedula string) (int, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.intentos[cedula], nil
}

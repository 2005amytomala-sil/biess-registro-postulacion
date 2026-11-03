package redis

import (
	"context"
	"sync"
	"time"
)

// InMemoryIntentosLimiter es una implementación en memoria para tests
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

func (l *InMemoryIntentosLimiter) ObtenerIntentos(ctx context.Context, cedula string) (int, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.intentos[cedula], nil
}

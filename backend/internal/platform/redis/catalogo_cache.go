package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	catalogoPaisesKey      = "catalogo:paises"
	catalogoProvinciasKey  = "catalogo:provincias"
	catalogoCantonesPrefix = "catalogo:cantones:"
)

type CatalogoCache struct {
	client *redis.Client
}

func NewCatalogoCache(redis *Redis) *CatalogoCache {
	return &CatalogoCache{
		client: redis.Client(),
	}
}

func (c *CatalogoCache) ObtenerPaises(ctx context.Context) ([]domain.Pais, error) {
	data, err := c.client.Get(ctx, catalogoPaisesKey).Result()
	if err != nil {
		return nil, err
	}

	var paises []domain.Pais
	if err := json.Unmarshal([]byte(data), &paises); err != nil {
		return nil, fmt.Errorf("error al deserializar paises: %w", err)
	}

	return paises, nil
}

func (c *CatalogoCache) GuardarPaises(ctx context.Context, paises []domain.Pais, ttl int) error {
	data, err := json.Marshal(paises)
	if err != nil {
		return fmt.Errorf("error al serializar paises: %w", err)
	}

	expiration := time.Duration(ttl) * time.Hour
	return c.client.Set(ctx, catalogoPaisesKey, data, expiration).Err()
}

func (c *CatalogoCache) ObtenerProvincias(ctx context.Context) ([]domain.Provincia, error) {
	data, err := c.client.Get(ctx, catalogoProvinciasKey).Result()
	if err != nil {
		return nil, err
	}

	var provincias []domain.Provincia
	if err := json.Unmarshal([]byte(data), &provincias); err != nil {
		return nil, fmt.Errorf("error al deserializar provincias: %w", err)
	}

	return provincias, nil
}

func (c *CatalogoCache) GuardarProvincias(ctx context.Context, provincias []domain.Provincia, ttl int) error {
	data, err := json.Marshal(provincias)
	if err != nil {
		return fmt.Errorf("error al serializar provincias: %w", err)
	}

	expiration := time.Duration(ttl) * time.Hour
	return c.client.Set(ctx, catalogoProvinciasKey, data, expiration).Err()
}

func (c *CatalogoCache) ObtenerCantonesPorProvincia(ctx context.Context, provinciaID int) ([]domain.Canton, error) {
	key := fmt.Sprintf("%s%d", catalogoCantonesPrefix, provinciaID)
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var cantones []domain.Canton
	if err := json.Unmarshal([]byte(data), &cantones); err != nil {
		return nil, fmt.Errorf("error al deserializar cantones: %w", err)
	}

	return cantones, nil
}

func (c *CatalogoCache) GuardarCantonesPorProvincia(ctx context.Context, provinciaID int, cantones []domain.Canton, ttl int) error {
	key := fmt.Sprintf("%s%d", catalogoCantonesPrefix, provinciaID)
	data, err := json.Marshal(cantones)
	if err != nil {
		return fmt.Errorf("error al serializar cantones: %w", err)
	}

	expiration := time.Duration(ttl) * time.Hour
	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *CatalogoCache) InvalidarCatalogos(ctx context.Context) error {
	// Invalidar claves específicas
	keys := []string{
		catalogoPaisesKey,
		catalogoProvinciasKey,
	}

	// Escanear y eliminar claves de cantones
	iter := c.client.Scan(ctx, 0, catalogoCantonesPrefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		log.Printf("error al escanear claves de cantones: %v", err)
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

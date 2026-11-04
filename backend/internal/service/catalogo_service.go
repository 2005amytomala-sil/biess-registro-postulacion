package service

import (
	"context"
	"fmt"
	"log"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
)

const (
	catalogoCacheTTL = 24 // 24 horas
)

type CatalogoService struct {
	repo   domain.CatalogoRepository
	cache  domain.CatalogoCache
	logger *log.Logger
}

func NewCatalogoService(repo domain.CatalogoRepository, cache domain.CatalogoCache, logger *log.Logger) *CatalogoService {
	return &CatalogoService{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

func (s *CatalogoService) ObtenerPaises(ctx context.Context) ([]domain.Pais, error) {
	// Intentar obtener del caché
	paises, err := s.cache.ObtenerPaises(ctx)
	if err == nil {
		return paises, nil
	}

	// Redis caído o cache miss: registrar y degradar a DB
	s.logger.Printf("cache miss/redis error paises: %v", err)

	// Consultar base de datos
	paises, err = s.repo.ObtenerPaises(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al obtener paises del repositorio: %w", err)
	}

	// Intentar guardar en caché (en background, no bloquear la respuesta)
	go func() {
		if err := s.cache.GuardarPaises(context.Background(), paises, catalogoCacheTTL); err != nil {
			s.logger.Printf("error al guardar paises en caché: %v", err)
		}
	}()

	return paises, nil
}

func (s *CatalogoService) ObtenerProvincias(ctx context.Context) ([]domain.Provincia, error) {
	// Intentar obtener del caché
	provincias, err := s.cache.ObtenerProvincias(ctx)
	if err == nil {
		return provincias, nil
	}

	// Redis caído o cache miss: registrar y degradar a DB
	s.logger.Printf("cache miss/redis error provincias: %v", err)

	// Consultar base de datos
	provincias, err = s.repo.ObtenerProvincias(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al obtener provincias del repositorio: %w", err)
	}

	// Intentar guardar en caché
	go func() {
		if err := s.cache.GuardarProvincias(context.Background(), provincias, catalogoCacheTTL); err != nil {
			s.logger.Printf("error al guardar provincias en caché: %v", err)
		}
	}()

	return provincias, nil
}

func (s *CatalogoService) ObtenerCantonesPorProvincia(ctx context.Context, provinciaID int) ([]domain.Canton, error) {
	// Verificar que la provincia existe
	existe, err := s.repo.ExisteProvincia(ctx, provinciaID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar existencia de provincia: %w", err)
	}
	if !existe {
		return nil, fmt.Errorf("provincia con ID %d no encontrada", provinciaID)
	}

	// Intentar obtener del caché
	cantones, err := s.cache.ObtenerCantonesPorProvincia(ctx, provinciaID)
	if err == nil {
		return cantones, nil
	}

	// Redis caído o cache miss: registrar y degradar a DB
	s.logger.Printf("cache miss/redis error cantones provincia %d: %v", provinciaID, err)

	// Consultar base de datos
	cantones, err = s.repo.ObtenerCantonesPorProvincia(ctx, provinciaID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener cantones del repositorio: %w", err)
	}

	// Intentar guardar en caché
	go func() {
		if err := s.cache.GuardarCantonesPorProvincia(context.Background(), provinciaID, cantones, catalogoCacheTTL); err != nil {
			s.logger.Printf("error al guardar cantones en caché: %v", err)
		}
	}()

	return cantones, nil
}

func (s *CatalogoService) InvalidarCache(ctx context.Context) error {
	s.logger.Println("invalidando caché de catálogos")
	return s.cache.InvalidarCatalogos(ctx)
}

package service

import (
	"bytes"
	"context"
	"log"
	"testing"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations for testing
type MockCatalogoRepository struct {
	mock.Mock
}

func (m *MockCatalogoRepository) ObtenerPaises(ctx context.Context) ([]domain.Pais, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Pais), args.Error(1)
}

func (m *MockCatalogoRepository) ObtenerProvincias(ctx context.Context) ([]domain.Provincia, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Provincia), args.Error(1)
}

func (m *MockCatalogoRepository) ObtenerCantonesPorProvincia(ctx context.Context, provinciaID int) ([]domain.Canton, error) {
	args := m.Called(ctx, provinciaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Canton), args.Error(1)
}

func (m *MockCatalogoRepository) ExisteProvincia(ctx context.Context, provinciaID int) (bool, error) {
	args := m.Called(ctx, provinciaID)
	return args.Bool(0), args.Error(1)
}

type MockCatalogoCache struct {
	mock.Mock
}

func (m *MockCatalogoCache) ObtenerPaises(ctx context.Context) ([]domain.Pais, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Pais), args.Error(1)
}

func (m *MockCatalogoCache) GuardarPaises(ctx context.Context, paises []domain.Pais, ttl int) error {
	args := m.Called(ctx, paises, ttl)
	return args.Error(0)
}

func (m *MockCatalogoCache) ObtenerProvincias(ctx context.Context) ([]domain.Provincia, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Provincia), args.Error(1)
}

func (m *MockCatalogoCache) GuardarProvincias(ctx context.Context, provincias []domain.Provincia, ttl int) error {
	args := m.Called(ctx, provincias, ttl)
	return args.Error(0)
}

func (m *MockCatalogoCache) ObtenerCantonesPorProvincia(ctx context.Context, provinciaID int) ([]domain.Canton, error) {
	args := m.Called(ctx, provinciaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Canton), args.Error(1)
}

func (m *MockCatalogoCache) GuardarCantonesPorProvincia(ctx context.Context, provinciaID int, cantones []domain.Canton, ttl int) error {
	args := m.Called(ctx, provinciaID, cantones, ttl)
	return args.Error(0)
}

func (m *MockCatalogoCache) InvalidarCatalogos(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestObtenerPaises_CacheHit(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCatalogoRepository{}
	mockCache := &MockCatalogoCache{}
	logger := log.New(&bytes.Buffer{}, "", 0)

	service := NewCatalogoService(mockRepo, mockCache, logger)

	paises := []domain.Pais{
		{ID: 1, Nombre: "Ecuador"},
		{ID: 2, Nombre: "Colombia"},
	}

	// Cache hit: mock cache retorna datos
	mockCache.On("ObtenerPaises", ctx).Return(paises, nil)

	result, err := service.ObtenerPaises(ctx)

	assert.NoError(t, err)
	assert.Equal(t, paises, result)
	mockCache.AssertExpectations(t)
}

func TestObtenerPaises_CacheMiss(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCatalogoRepository{}
	mockCache := &MockCatalogoCache{}
	logger := log.New(&bytes.Buffer{}, "", 0)

	service := NewCatalogoService(mockRepo, mockCache, logger)

	paises := []domain.Pais{
		{ID: 1, Nombre: "Ecuador"},
		{ID: 2, Nombre: "Colombia"},
	}

	// Cache miss: cache retorna error, repository retorna datos
	mockCache.On("ObtenerPaises", ctx).Return(nil, assert.AnError)
	mockRepo.On("ObtenerPaises", ctx).Return(paises, nil)
	// No verificamos GuardarPaises porque se ejecuta en background

	result, err := service.ObtenerPaises(ctx)

	assert.NoError(t, err)
	assert.Equal(t, paises, result)
	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestObtenerPaises_RedisCaido(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCatalogoRepository{}
	mockCache := &MockCatalogoCache{}
	logger := log.New(&bytes.Buffer{}, "", 0)

	service := NewCatalogoService(mockRepo, mockCache, logger)

	paises := []domain.Pais{
		{ID: 1, Nombre: "Ecuador"},
	}

	// Redis caído: cache retorna error, degrada a DB
	mockCache.On("ObtenerPaises", ctx).Return(nil, assert.AnError)
	mockRepo.On("ObtenerPaises", ctx).Return(paises, nil)

	result, err := service.ObtenerPaises(ctx)

	// Debe retornar datos de DB a pesar de error de Redis
	assert.NoError(t, err)
	assert.Equal(t, paises, result)
	mockRepo.AssertExpectations(t)
}

func TestObtenerProvincias_CacheHit(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCatalogoRepository{}
	mockCache := &MockCatalogoCache{}
	logger := log.New(&bytes.Buffer{}, "", 0)

	service := NewCatalogoService(mockRepo, mockCache, logger)

	provincias := []domain.Provincia{
		{ID: 1, Nombre: "Pichincha"},
		{ID: 2, Nombre: "Guayas"},
	}

	mockCache.On("ObtenerProvincias", ctx).Return(provincias, nil)

	result, err := service.ObtenerProvincias(ctx)

	assert.NoError(t, err)
	assert.Equal(t, provincias, result)
}

func TestObtenerProvincias_CacheMiss(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCatalogoRepository{}
	mockCache := &MockCatalogoCache{}
	logger := log.New(&bytes.Buffer{}, "", 0)

	service := NewCatalogoService(mockRepo, mockCache, logger)

	provincias := []domain.Provincia{
		{ID: 1, Nombre: "Pichincha"},
	}

	mockCache.On("ObtenerProvincias", ctx).Return(nil, assert.AnError)
	mockRepo.On("ObtenerProvincias", ctx).Return(provincias, nil)

	result, err := service.ObtenerProvincias(ctx)

	assert.NoError(t, err)
	assert.Equal(t, provincias, result)
}

func TestObtenerCantonesPorProvincia_CacheHit(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCatalogoRepository{}
	mockCache := &MockCatalogoCache{}
	logger := log.New(&bytes.Buffer{}, "", 0)

	service := NewCatalogoService(mockRepo, mockCache, logger)

	cantones := []domain.Canton{
		{ID: 1, ProvinciaID: 19, Nombre: "Quito"},
		{ID: 2, ProvinciaID: 19, Nombre: "Cayambe"},
	}

	mockRepo.On("ExisteProvincia", ctx, 19).Return(true, nil)
	mockCache.On("ObtenerCantonesPorProvincia", ctx, 19).Return(cantones, nil)

	result, err := service.ObtenerCantonesPorProvincia(ctx, 19)

	assert.NoError(t, err)
	assert.Equal(t, cantones, result)
	mockCache.AssertExpectations(t)
}

func TestObtenerCantonesPorProvincia_CacheMiss(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCatalogoRepository{}
	mockCache := &MockCatalogoCache{}
	logger := log.New(&bytes.Buffer{}, "", 0)

	service := NewCatalogoService(mockRepo, mockCache, logger)

	cantones := []domain.Canton{
		{ID: 1, ProvinciaID: 19, Nombre: "Quito"},
	}

	mockRepo.On("ExisteProvincia", ctx, 19).Return(true, nil)
	mockCache.On("ObtenerCantonesPorProvincia", ctx, 19).Return(nil, assert.AnError)
	mockRepo.On("ObtenerCantonesPorProvincia", ctx, 19).Return(cantones, nil)

	result, err := service.ObtenerCantonesPorProvincia(ctx, 19)

	assert.NoError(t, err)
	assert.Equal(t, cantones, result)
}

func TestObtenerCantonesPorProvincia_ProvinciaNoExiste(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCatalogoRepository{}
	mockCache := &MockCatalogoCache{}
	logger := log.New(&bytes.Buffer{}, "", 0)

	service := NewCatalogoService(mockRepo, mockCache, logger)

	mockRepo.On("ExisteProvincia", ctx, 999).Return(false, nil)

	_, err := service.ObtenerCantonesPorProvincia(ctx, 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no encontrada")
}

func TestObtenerCantonesPorProvincia_RedisCaido(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCatalogoRepository{}
	mockCache := &MockCatalogoCache{}
	logger := log.New(&bytes.Buffer{}, "", 0)

	service := NewCatalogoService(mockRepo, mockCache, logger)

	cantones := []domain.Canton{
		{ID: 1, ProvinciaID: 19, Nombre: "Quito"},
	}

	mockRepo.On("ExisteProvincia", ctx, 19).Return(true, nil)
	mockCache.On("ObtenerCantonesPorProvincia", ctx, 19).Return(nil, assert.AnError)
	mockRepo.On("ObtenerCantonesPorProvincia", ctx, 19).Return(cantones, nil)

	result, err := service.ObtenerCantonesPorProvincia(ctx, 19)

	// Debe retornar datos de DB a pesar de error de Redis
	assert.NoError(t, err)
	assert.Equal(t, cantones, result)
}

func TestInvalidarCache(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCatalogoRepository{}
	mockCache := &MockCatalogoCache{}
	logger := log.New(&bytes.Buffer{}, "", 0)

	service := NewCatalogoService(mockRepo, mockCache, logger)

	mockCache.On("InvalidarCatalogos", ctx).Return(nil)

	err := service.InvalidarCache(ctx)

	assert.NoError(t, err)
	mockCache.AssertExpectations(t)
}

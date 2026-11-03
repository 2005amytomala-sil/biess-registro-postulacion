package service

import (
	"context"
	"testing"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/jwt"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Mock implementations for testing
type MockPostulanteRepository struct {
	mock.Mock
}

func (m *MockPostulanteRepository) Create(ctx context.Context, postulante *domain.Postulante) error {
	args := m.Called(ctx, postulante)
	return args.Error(0)
}

func (m *MockPostulanteRepository) FindByCedula(ctx context.Context, cedula string) (*domain.Postulante, error) {
	args := m.Called(ctx, cedula)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Postulante), args.Error(1)
}

func (m *MockPostulanteRepository) FindByCorreo(ctx context.Context, correo string) (*domain.Postulante, error) {
	args := m.Called(ctx, correo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Postulante), args.Error(1)
}

func (m *MockPostulanteRepository) ExistsByCedula(ctx context.Context, cedula string) (bool, error) {
	args := m.Called(ctx, cedula)
	return args.Bool(0), args.Error(1)
}

func (m *MockPostulanteRepository) ExistsByCorreo(ctx context.Context, correo string) (bool, error) {
	args := m.Called(ctx, correo)
	return args.Bool(0), args.Error(1)
}

func (m *MockPostulanteRepository) UpdatePassword(ctx context.Context, id int, passwordHash string) error {
	args := m.Called(ctx, id, passwordHash)
	return args.Error(0)
}

type MockIntentosLimiter struct {
	mock.Mock
}

func (m *MockIntentosLimiter) IncrementarIntento(ctx context.Context, cedula string) (int, error) {
	args := m.Called(ctx, cedula)
	return args.Int(0), args.Error(1)
}

func (m *MockIntentosLimiter) VerificarBloqueo(ctx context.Context, cedula string) (bool, int, error) {
	args := m.Called(ctx, cedula)
	return args.Bool(0), args.Int(1), args.Error(2)
}

func (m *MockIntentosLimiter) ReiniciarIntentos(ctx context.Context, cedula string) error {
	args := m.Called(ctx, cedula)
	return args.Error(0)
}

func (m *MockIntentosLimiter) EstablecerBloqueo(ctx context.Context, cedula string, ttl int) error {
	args := m.Called(ctx, cedula, ttl)
	return args.Error(0)
}

type MockTokenStorage struct {
	mock.Mock
}

func (m *MockTokenStorage) GuardarToken(ctx context.Context, cedula string, token string) error {
	args := m.Called(ctx, cedula, token)
	return args.Error(0)
}

func (m *MockTokenStorage) ValidarToken(ctx context.Context, cedula string, token string) (bool, error) {
	args := m.Called(ctx, cedula, token)
	return args.Bool(0), args.Error(1)
}

func (m *MockTokenStorage) EliminarToken(ctx context.Context, cedula string) error {
	args := m.Called(ctx, cedula)
	return args.Error(0)
}

func TestVerificarIdentidad_CasoFeliz(t *testing.T) {
	ctx := context.Background()
	registroCivil := mocks.NewRegistroCivilMock()
	intentosLimiter := &MockIntentosLimiter{}
	tokenManager := jwt.NewTokenManager("test-secret", 24)
	tokenStorage := &MockTokenStorage{}
	mockRepo := &MockPostulanteRepository{}

	authService := NewAuthService(mockRepo, registroCivil, intentosLimiter, tokenManager, tokenStorage)

	// Test con datos válidos
	data, err := authService.VerificarIdentidad(ctx, "0101010101", "01/01/2010", "")

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "0101010101", data.Cedula)
	assert.Equal(t, "Ecuatoriana", data.Nacionalidad)
}

func TestVerificarIdentidad_IdentidadInvalida(t *testing.T) {
	ctx := context.Background()
	registroCivil := mocks.NewRegistroCivilMock()
	intentosLimiter := &MockIntentosLimiter{}
	tokenManager := jwt.NewTokenManager("test-secret", 24)
	tokenStorage := &MockTokenStorage{}
	mockRepo := &MockPostulanteRepository{}

	authService := NewAuthService(mockRepo, registroCivil, intentosLimiter, tokenManager, tokenStorage)

	// Test con cédula no existente
	_, err := authService.VerificarIdentidad(ctx, "9999999999", "01/01/2010", "")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrIdentidadInvalida, err)
}

func TestVerificarIdentidad_NacionalidadNoEcuatoriana(t *testing.T) {
	ctx := context.Background()
	registroCivil := mocks.NewRegistroCivilMock()
	intentosLimiter := &MockIntentosLimiter{}
	tokenManager := jwt.NewTokenManager("test-secret", 24)
	tokenStorage := &MockTokenStorage{}
	mockRepo := &MockPostulanteRepository{}

	authService := NewAuthService(mockRepo, registroCivil, intentosLimiter, tokenManager, tokenStorage)

	// Test con nacionalidad no ecuatoriana
	_, err := authService.VerificarIdentidad(ctx, "0303030303", "20/03/2015", "")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrIdentidadInvalida, err)
}

func TestLogin_CasoFeliz(t *testing.T) {
	ctx := context.Background()
	intentosLimiter := &MockIntentosLimiter{}
	tokenManager := jwt.NewTokenManager("test-secret", 24)
	tokenStorage := &MockTokenStorage{}
	mockRepo := &MockPostulanteRepository{}

	authService := NewAuthService(mockRepo, mocks.NewRegistroCivilMock(), intentosLimiter, tokenManager, tokenStorage)

	// Generar hash real de "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	intentosLimiter.On("VerificarBloqueo", ctx, "0101010101").Return(false, 0, nil)
	intentosLimiter.On("IncrementarIntento", ctx, "0101010101").Return(1, nil)
	mockRepo.On("FindByCedula", ctx, "0101010101").Return(&domain.Postulante{
		ID: 1, Cedula: "0101010101", Correo: "test@example.com",
		PasswordHash: string(hashedPassword),
	}, nil)
	intentosLimiter.On("ReiniciarIntentos", ctx, "0101010101").Return(nil)

	response, err := authService.Login(ctx, "0101010101", "password")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Token)
}

func TestLogin_CredencialesInvalidas(t *testing.T) {
	ctx := context.Background()
	intentosLimiter := &MockIntentosLimiter{}
	tokenManager := jwt.NewTokenManager("test-secret", 24)
	tokenStorage := &MockTokenStorage{}
	mockRepo := &MockPostulanteRepository{}

	authService := NewAuthService(mockRepo, mocks.NewRegistroCivilMock(), intentosLimiter, tokenManager, tokenStorage)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	intentosLimiter.On("VerificarBloqueo", ctx, "0101010101").Return(false, 0, nil)
	intentosLimiter.On("IncrementarIntento", ctx, "0101010101").Return(1, nil)
	mockRepo.On("FindByCedula", ctx, "0101010101").Return(&domain.Postulante{
		ID: 1, Cedula: "0101010101", Correo: "test@example.com",
		PasswordHash: string(hashedPassword),
	}, nil)

	_, err := authService.Login(ctx, "0101010101", "wrongpassword")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrCredencialesInvalidas, err)
}

func TestLogin_Bloqueado(t *testing.T) {
	ctx := context.Background()
	intentosLimiter := &MockIntentosLimiter{}
	tokenManager := jwt.NewTokenManager("test-secret", 24)
	tokenStorage := &MockTokenStorage{}
	mockRepo := &MockPostulanteRepository{}

	authService := NewAuthService(mockRepo, mocks.NewRegistroCivilMock(), intentosLimiter, tokenManager, tokenStorage)
	intentosLimiter.On("VerificarBloqueo", ctx, "0101010101").Return(true, 600, nil)

	_, err := authService.Login(ctx, "0101010101", "password")
	assert.Error(t, err)
}

func TestRegistrar_CasoFeliz(t *testing.T) {
	ctx := context.Background()
	intentosLimiter := &MockIntentosLimiter{}
	tokenManager := jwt.NewTokenManager("test-secret", 24)
	tokenStorage := &MockTokenStorage{}
	mockRepo := &MockPostulanteRepository{}

	authService := NewAuthService(mockRepo, mocks.NewRegistroCivilMock(), intentosLimiter, tokenManager, tokenStorage)
	mockRepo.On("ExistsByCedula", ctx, "0101010101").Return(false, nil)
	mockRepo.On("ExistsByCorreo", ctx, "new@example.com").Return(false, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Postulante")).Return(nil).Run(func(args mock.Arguments) {
		postulante := args.Get(1).(*domain.Postulante)
		postulante.ID = 1
	})

	response, err := authService.Registrar(ctx, "0101010101", "new@example.com", "password123", "password123")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Token)
}

func TestRegistrar_YaRegistrado(t *testing.T) {
	ctx := context.Background()
	intentosLimiter := &MockIntentosLimiter{}
	tokenManager := jwt.NewTokenManager("test-secret", 24)
	tokenStorage := &MockTokenStorage{}
	mockRepo := &MockPostulanteRepository{}

	authService := NewAuthService(mockRepo, mocks.NewRegistroCivilMock(), intentosLimiter, tokenManager, tokenStorage)
	mockRepo.On("ExistsByCedula", ctx, "0101010101").Return(true, nil)

	_, err := authService.Registrar(ctx, "0101010101", "new@example.com", "password123", "password123")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrYaRegistrado, err)
}

func TestRegistrar_PasswordMismatch(t *testing.T) {
	ctx := context.Background()
	intentosLimiter := &MockIntentosLimiter{}
	tokenManager := jwt.NewTokenManager("test-secret", 24)
	tokenStorage := &MockTokenStorage{}
	mockRepo := &MockPostulanteRepository{}

	authService := NewAuthService(mockRepo, mocks.NewRegistroCivilMock(), intentosLimiter, tokenManager, tokenStorage)

	_, err := authService.Registrar(ctx, "0101010101", "new@example.com", "password123", "different123")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrPasswordMismatch, err)
}

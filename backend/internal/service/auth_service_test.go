package service

import (
	"context"
	"errors"
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
	registroLimiter := &MockIntentosLimiter{}
	loginLimiter := &MockIntentosLimiter{}
	tokenManager := jwt.NewTokenManager("test-secret", 24)
	tokenStorage := &MockTokenStorage{}
	mockRepo := &MockPostulanteRepository{}

	authService := NewAuthService(mockRepo, registroCivil, registroLimiter, loginLimiter, tokenManager, tokenStorage)

	// Configurar expectativas del mock
	registroLimiter.On("VerificarBloqueo", mock.Anything, mock.Anything).Return(false, 0, nil)
	registroLimiter.On("ReiniciarIntentos", mock.Anything, mock.Anything).Return(nil)

	// Test con datos válidos
	data, err := authService.VerificarIdentidad(ctx, "0101010101", "01/01/2010", "")

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "0101010101", data.Cedula)
	assert.Equal(t, "Ecuatoriana", data.Nacionalidad)

	registroLimiter.AssertExpectations(t)
}

func TestVerificarIdentidad_IdentidadInvalida(t *testing.T) {
	ctx := context.Background()
	registroCivil := mocks.NewRegistroCivilMock()
	registroLimiter := &MockIntentosLimiter{}
	loginLimiter := &MockIntentosLimiter{}
	tokenManager := jwt.NewTokenManager("test-secret", 24)
	tokenStorage := &MockTokenStorage{}
	mockRepo := &MockPostulanteRepository{}

	authService := NewAuthService(mockRepo, registroCivil, registroLimiter, loginLimiter, tokenManager, tokenStorage)

	// Configurar expectativas del mock
	registroLimiter.On("VerificarBloqueo", mock.Anything, mock.Anything).Return(false, 0, nil)
	registroLimiter.On("IncrementarIntento", mock.Anything, mock.Anything).Return(1, nil)

	// Test con cédula existente pero fecha incorrecta (simula identidad inválida)
	_, err := authService.VerificarIdentidad(ctx, "0101010101", "01/01/2000", "")

	assert.Error(t, err)
	assert.Equal(t, domain.ErrIdentidadInvalida, err)

	registroLimiter.AssertExpectations(t)
}

type authTestDeps struct {
	repo            *MockPostulanteRepository
	registroLimiter *MockIntentosLimiter
	loginLimiter    *MockIntentosLimiter
	service         *AuthService
}

func newAuthTestDeps() *authTestDeps {
	d := &authTestDeps{
		repo:            &MockPostulanteRepository{},
		registroLimiter: &MockIntentosLimiter{},
		loginLimiter:    &MockIntentosLimiter{},
	}
	d.service = NewAuthService(d.repo, mocks.NewRegistroCivilMock(), d.registroLimiter, d.loginLimiter, jwt.NewTokenManager("test-secret", 24), &MockTokenStorage{})
	return d
}

func TestVerificarIdentidad_NacionalidadNoEcuatoriana(t *testing.T) {
	d := newAuthTestDeps()
	d.registroLimiter.On("VerificarBloqueo", mock.Anything, mock.Anything).Return(false, 0, nil)

	_, err := d.service.VerificarIdentidad(context.Background(), "0303030303", "20/03/2015", "")

	assert.Equal(t, domain.ErrIdentidadInvalida, err)
	d.registroLimiter.AssertExpectations(t)
	d.loginLimiter.AssertNotCalled(t, "VerificarBloqueo", mock.Anything, mock.Anything)
}

func TestVerificarIdentidad_SinFechaNiCodigo(t *testing.T) {
	d := newAuthTestDeps()
	d.registroLimiter.On("VerificarBloqueo", mock.Anything, mock.Anything).Return(false, 0, nil)
	d.registroLimiter.On("IncrementarIntento", mock.Anything, "0101010101").Return(1, nil)

	_, err := d.service.VerificarIdentidad(context.Background(), "0101010101", "", "")

	assert.Equal(t, domain.ErrIdentidadInvalida, err)
	d.registroLimiter.AssertExpectations(t)
}

func TestVerificarIdentidad_TercerFalloBloquea(t *testing.T) {
	d := newAuthTestDeps()
	d.registroLimiter.On("VerificarBloqueo", mock.Anything, mock.Anything).Return(false, 0, nil)
	d.registroLimiter.On("IncrementarIntento", mock.Anything, "0101010101").Return(3, nil)
	d.registroLimiter.On("EstablecerBloqueo", mock.Anything, "0101010101", bloqueoTTLSeconds).Return(nil)

	_, err := d.service.VerificarIdentidad(context.Background(), "0101010101", "01/01/2000", "")

	var appErr *domain.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, domain.ErrBloqueado.Code, appErr.Code)
	assert.Equal(t, "Registro fallido, inténtelo nuevamente en 10 minutos", appErr.Message)
	assert.Equal(t, bloqueoTTLSeconds, appErr.RetryAfter)
	d.registroLimiter.AssertExpectations(t)
}

func TestVerificarIdentidad_ErrorEstableciendoBloqueo(t *testing.T) {
	d := newAuthTestDeps()
	redisErr := errors.New("redis caído")
	d.registroLimiter.On("VerificarBloqueo", mock.Anything, mock.Anything).Return(false, 0, nil)
	d.registroLimiter.On("IncrementarIntento", mock.Anything, mock.Anything).Return(3, nil)
	d.registroLimiter.On("EstablecerBloqueo", mock.Anything, mock.Anything, mock.Anything).Return(redisErr)

	_, err := d.service.VerificarIdentidad(context.Background(), "0101010101", "01/01/2000", "")

	assert.ErrorIs(t, err, redisErr)
}

func TestVerificarIdentidad_BloqueadoDevuelveTTLReal(t *testing.T) {
	d := newAuthTestDeps()
	d.registroLimiter.On("VerificarBloqueo", mock.Anything, "0101010101").Return(true, 437, nil)

	_, err := d.service.VerificarIdentidad(context.Background(), "0101010101", "01/01/2010", "")

	var appErr *domain.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, domain.ErrBloqueado.Code, appErr.Code)
	assert.Equal(t, 437, appErr.RetryAfter)
}

func TestLogin_CasoFeliz(t *testing.T) {
	d := newAuthTestDeps()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	d.loginLimiter.On("VerificarBloqueo", mock.Anything, mock.Anything).Return(false, 0, nil)
	d.repo.On("FindByCedula", mock.Anything, mock.Anything).Return(&domain.Postulante{
		ID: 1, Cedula: "0101010101", Correo: "test@example.com",
		PasswordHash: string(hashedPassword),
	}, nil)
	d.loginLimiter.On("ReiniciarIntentos", mock.Anything, mock.Anything).Return(nil)

	response, err := d.service.Login(context.Background(), "0101010101", "password")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Token)

	d.loginLimiter.AssertExpectations(t)
	d.registroLimiter.AssertNotCalled(t, "VerificarBloqueo", mock.Anything, mock.Anything)
}

func TestLogin_CredencialesInvalidas(t *testing.T) {
	d := newAuthTestDeps()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	d.loginLimiter.On("VerificarBloqueo", mock.Anything, mock.AnythingOfType("string")).Return(false, 0, nil)
	d.loginLimiter.On("IncrementarIntento", mock.Anything, "0101010101").Return(1, nil)
	d.repo.On("FindByCedula", mock.Anything, mock.AnythingOfType("string")).Return(&domain.Postulante{
		ID: 1, Cedula: "0101010101", Correo: "test@example.com",
		PasswordHash: string(hashedPassword),
	}, nil)

	_, err := d.service.Login(context.Background(), "0101010101", "wrongpassword")
	assert.Equal(t, domain.ErrCredencialesInvalidas, err)
	assert.Equal(t, "Sus credenciales de ingreso son incorrectas", domain.ErrCredencialesInvalidas.Message)

	d.loginLimiter.AssertExpectations(t)
	d.registroLimiter.AssertNotCalled(t, "IncrementarIntento", mock.Anything, mock.Anything)
}

func TestLogin_TercerFalloBloquea(t *testing.T) {
	d := newAuthTestDeps()
	d.loginLimiter.On("VerificarBloqueo", mock.Anything, mock.Anything).Return(false, 0, nil)
	d.loginLimiter.On("IncrementarIntento", mock.Anything, "0101010101").Return(3, nil)
	d.loginLimiter.On("EstablecerBloqueo", mock.Anything, "0101010101", bloqueoTTLSeconds).Return(nil)
	d.repo.On("FindByCedula", mock.Anything, mock.Anything).Return(nil, domain.ErrNotFound)

	_, err := d.service.Login(context.Background(), "0101010101", "password")

	var appErr *domain.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, domain.ErrBloqueado.Code, appErr.Code)
	assert.Equal(t, bloqueoTTLSeconds, appErr.RetryAfter)
	d.loginLimiter.AssertExpectations(t)
}

func TestLogin_Bloqueado(t *testing.T) {
	ctx := context.Background()
	d := newAuthTestDeps()
	d.loginLimiter.On("VerificarBloqueo", ctx, "0101010101").Return(true, 321, nil)

	_, err := d.service.Login(ctx, "0101010101", "password")

	var appErr *domain.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, domain.ErrBloqueado.Code, appErr.Code)
	assert.Equal(t, 321, appErr.RetryAfter)
}

func TestRegistrar_CasoFeliz(t *testing.T) {
	ctx := context.Background()
	d := newAuthTestDeps()
	d.registroLimiter.On("VerificarBloqueo", ctx, "0101010101").Return(false, 0, nil)
	d.registroLimiter.On("ReiniciarIntentos", ctx, "0101010101").Return(nil)
	d.repo.On("ExistsByCedula", ctx, "0101010101").Return(false, nil)
	d.repo.On("ExistsByCorreo", ctx, "new@example.com").Return(false, nil)
	d.repo.On("Create", ctx, mock.AnythingOfType("*domain.Postulante")).Return(nil).Run(func(args mock.Arguments) {
		args.Get(1).(*domain.Postulante).ID = 1
	})

	response, err := d.service.Registrar(ctx, "0101010101", "", "1234567890", "new@example.com", "password123", "password123")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Token)
	d.registroLimiter.AssertExpectations(t)
	d.repo.AssertExpectations(t)
}

func TestRegistrar_SinDatosIdentidad(t *testing.T) {
	d := newAuthTestDeps()

	_, err := d.service.Registrar(context.Background(), "0101010101", "", "", "new@example.com", "password123", "password123")

	assert.Equal(t, domain.ErrIdentidadInvalida, err)
	d.repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestRegistrar_IdentidadIncorrecta(t *testing.T) {
	d := newAuthTestDeps()
	d.registroLimiter.On("VerificarBloqueo", mock.Anything, mock.Anything).Return(false, 0, nil)
	d.registroLimiter.On("IncrementarIntento", mock.Anything, "0101010101").Return(1, nil)

	_, err := d.service.Registrar(context.Background(), "0101010101", "01/01/2000", "", "new@example.com", "password123", "password123")

	assert.Equal(t, domain.ErrIdentidadInvalida, err)
	assert.Equal(t, "La información ingresada es incorrecta", domain.ErrIdentidadInvalida.Message)
	d.registroLimiter.AssertExpectations(t)
	d.repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestRegistrar_NacionalidadNoEcuatoriana(t *testing.T) {
	d := newAuthTestDeps()
	d.registroLimiter.On("VerificarBloqueo", mock.Anything, mock.Anything).Return(false, 0, nil)

	_, err := d.service.Registrar(context.Background(), "0303030303", "20/03/2015", "", "new@example.com", "password123", "password123")

	assert.Equal(t, domain.ErrIdentidadInvalida, err)
	d.repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestRegistrar_Bloqueado(t *testing.T) {
	d := newAuthTestDeps()
	d.registroLimiter.On("VerificarBloqueo", mock.Anything, "0101010101").Return(true, 512, nil)

	_, err := d.service.Registrar(context.Background(), "0101010101", "01/01/2010", "", "new@example.com", "password123", "password123")

	var appErr *domain.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, domain.ErrBloqueado.Code, appErr.Code)
	assert.Equal(t, 512, appErr.RetryAfter)
	d.repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestRegistrar_YaRegistrado(t *testing.T) {
	ctx := context.Background()
	d := newAuthTestDeps()
	d.registroLimiter.On("VerificarBloqueo", ctx, "0101010101").Return(false, 0, nil)
	d.registroLimiter.On("ReiniciarIntentos", ctx, "0101010101").Return(nil)
	d.repo.On("ExistsByCedula", ctx, "0101010101").Return(true, nil)

	_, err := d.service.Registrar(ctx, "0101010101", "01/01/2010", "", "new@example.com", "password123", "password123")
	assert.Equal(t, domain.ErrYaRegistrado, err)
}

func TestRegistrar_PasswordMismatch(t *testing.T) {
	d := newAuthTestDeps()

	_, err := d.service.Registrar(context.Background(), "0101010101", "01/01/2010", "", "new@example.com", "password123", "different123")
	assert.Equal(t, domain.ErrPasswordMismatch, err)
}

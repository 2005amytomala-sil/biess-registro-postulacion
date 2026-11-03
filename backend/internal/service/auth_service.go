package service

import (
	"context"
	"fmt"
	"log"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/redis"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	postulanteRepo  domain.PostulanteRepository
	registroCivil   domain.RegistroCivilClient
	intentosLimiter *redis.IntentosLimiter
	tokenManager    domain.TokenManager
	tokenStorage    *redis.TokenStorage
}

func NewAuthService(
	postulanteRepo domain.PostulanteRepository,
	registroCivil domain.RegistroCivilClient,
	intentosLimiter *redis.IntentosLimiter,
	tokenManager domain.TokenManager,
	tokenStorage *redis.TokenStorage,
) *AuthService {
	return &AuthService{
		postulanteRepo:  postulanteRepo,
		registroCivil:   registroCivil,
		intentosLimiter: intentosLimiter,
		tokenManager:    tokenManager,
		tokenStorage:    tokenStorage,
	}
}

// VerificarIdentidad verifica la identidad del postulante con el Registro Civil
func (s *AuthService) VerificarIdentidad(ctx context.Context, cedula string, fechaExpedicion string, codigoDactilar string) (*domain.RegistroCivilData, error) {
	// Verificar con el mock del Registro Civil
	data, err := s.registroCivil.VerificarIdentidad(ctx, cedula, fechaExpedicion, codigoDactilar)
	if err != nil {
		return nil, err
	}

	// Validar nacionalidad ecuatoriana
	if data.Nacionalidad != "Ecuatoriana" {
		return nil, domain.ErrIdentidadInvalida
	}

	return data, nil
}

// Registrar registra un nuevo postulante
func (s *AuthService) Registrar(ctx context.Context, cedula string, correo string, password string, confirmPassword string) (*domain.AuthResponse, error) {
	// Validar que las contraseñas coincidan
	if password != confirmPassword {
		return nil, domain.ErrPasswordMismatch
	}

	// Verificar si ya existe por cédula
	exists, err := s.postulanteRepo.ExistsByCedula(ctx, cedula)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrYaRegistrado
	}

	// Verificar si ya existe por correo
	exists, err = s.postulanteRepo.ExistsByCorreo(ctx, correo)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrYaRegistrado
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hasheando contraseña: %w", err)
	}

	// Crear postulante
	postulante := &domain.Postulante{
		Cedula:       cedula,
		Correo:       correo,
		PasswordHash: string(hashedPassword),
	}

	if err := s.postulanteRepo.Create(ctx, postulante); err != nil {
		return nil, err
	}

	// Generar token
	token, err := s.tokenManager.GenerarToken(postulante.ID, postulante.Cedula)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User:  *postulante,
	}, nil
}

// Login autentica a un postulante
func (s *AuthService) Login(ctx context.Context, cedula string, password string) (*domain.AuthResponse, error) {
	// Verificar bloqueo
	err := s.intentosLimiter.VerificarYBloquearSiNecesario(ctx, cedula)
	if err != nil {
		return nil, err
	}

	// Buscar postulante por cédula
	postulante, err := s.postulanteRepo.FindByCedula(ctx, cedula)
	if err != nil {
		// Incrementar intentos
		s.intentosLimiter.IncrementarIntento(ctx, cedula)
		return nil, domain.ErrCredencialesInvalidas
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(postulante.PasswordHash), []byte(password)); err != nil {
		// Incrementar intentos
		s.intentosLimiter.IncrementarIntento(ctx, cedula)
		return nil, domain.ErrCredencialesInvalidas
	}

	// Reiniciar intentos (login exitoso)
	s.intentosLimiter.ReiniciarIntentos(ctx, cedula)

	// Generar token
	token, err := s.tokenManager.GenerarToken(postulante.ID, postulante.Cedula)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User:  *postulante,
	}, nil
}

// Recuperar genera un token de recuperación y simula el envío de correo
func (s *AuthService) Recuperar(ctx context.Context, cedula string, fechaExpedicion string, codigoDactilar string) error {
	// Verificar identidad primero
	_, err := s.VerificarIdentidad(ctx, cedula, fechaExpedicion, codigoDactilar)
	if err != nil {
		return err
	}

	// Buscar postulante por cédula
	postulante, err := s.postulanteRepo.FindByCedula(ctx, cedula)
	if err != nil {
		return domain.ErrNotFound
	}

	// Generar token de recuperación (usando el mismo token manager)
	token, err := s.tokenManager.GenerarToken(postulante.ID, postulante.Cedula)
	if err != nil {
		return err
	}

	// Guardar token en Redis con TTL de 15 minutos
	if err := s.tokenStorage.GuardarToken(ctx, cedula, token); err != nil {
		return err
	}

	// Simular envío de correo con log
	log.Printf("[RECUPERACION] Token generado para %s (%s): %s", postulante.Cedula, postulante.Correo, token)

	return nil
}

// Restablecer permite restablecer la contraseña usando un token
func (s *AuthService) Restablecer(ctx context.Context, cedula string, token string, password string, confirmPassword string) error {
	// Validar que las contraseñas coincidan
	if password != confirmPassword {
		return domain.ErrPasswordMismatch
	}

	// Validar token contra Redis
	valid, err := s.tokenStorage.ValidarToken(ctx, cedula, token)
	if err != nil {
		return err
	}
	if !valid {
		return domain.ErrTokenInvalido
	}

	// Buscar postulante por cédula
	postulante, err := s.postulanteRepo.FindByCedula(ctx, cedula)
	if err != nil {
		return domain.ErrNotFound
	}

	// Hash de la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hasheando contraseña: %w", err)
	}

	// Actualizar contraseña en la base de datos
	if err := s.postulanteRepo.UpdatePassword(ctx, postulante.ID, string(hashedPassword)); err != nil {
		return err
	}

	// Eliminar token de Redis (usado)
	s.tokenStorage.EliminarToken(ctx, cedula)

	return nil
}

package domain

import (
	"fmt"
)

type AppError struct {
	Code       string
	Message    string
	Err        error
	RetryAfter int // Tiempo restante en segundos para bloqueos (opcional)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

var (
	ErrNotFound       = NewAppError("NOT_FOUND", "Recurso no encontrado", nil)
	ErrUnauthorized   = NewAppError("UNAUTHORIZED", "No autorizado", nil)
	ErrForbidden      = NewAppError("FORBIDDEN", "Acceso prohibido", nil)
	ErrBadRequest     = NewAppError("BAD_REQUEST", "Solicitud inválida", nil)
	ErrConflict       = NewAppError("CONFLICT", "Conflicto de datos", nil)
	ErrInternal       = NewAppError("INTERNAL", "Error interno del servidor", nil)
	ErrValidation      = NewAppError("VALIDATION", "Error de validación", nil)
	ErrTooManyRequests = NewAppError("TOO_MANY_REQUESTS", "Demasiadas solicitudes", nil)

	// Auth specific errors
	ErrIdentidadInvalida  = NewAppError("IDENTIDAD_INVALIDA", "La información ingresada es incorrecta", nil)
	ErrBloqueado          = NewAppError("BLOQUEADO", "Registro fallido, inténtelo nuevamente en 10 minutos", nil)
	ErrCredencialesInvalidas = NewAppError("CREDENCIALES_INVALIDAS", "Sus credenciales de ingreso son incorrectas", nil)
	ErrYaRegistrado       = NewAppError("YA_REGISTRADO", "El usuario ya está registrado", nil)
	ErrPasswordMismatch   = NewAppError("PASSWORD_MISMATCH", "Las contraseñas no coinciden", nil)
	ErrTokenInvalido      = NewAppError("TOKEN_INVALIDO", "Token inválido o expirado", nil)
	ErrNoAutorizado       = NewAppError("NO_AUTORIZADO", "Sesión inválida o expirada", nil)

	// Postulacion specific errors
	ErrNacionalidadNoPermitida      = NewAppError("NACIONALIDAD_NO_PERMITIDA", "Solo se permite nacionalidad ecuatoriana", nil)
	ErrCantonNoPerteneceAProvincia  = NewAppError("CANTON_NO_PERTENECE_PROVINCIA", "El cantón no pertenece a la provincia seleccionada", nil)
	ErrSectorNoValido               = NewAppError("SECTOR_NO_VALIDO", "Sector no válido para postulación", nil)
	ErrListaFamiliaresVacia         = NewAppError("LISTA_FAMILIARES_VACIA", "Debe agregar al menos un familiar", nil)
	ErrFechaFinAnteriorInicio       = NewAppError("FECHA_FIN_ANTERIOR_INICIO", "La fecha de fin no puede ser anterior a la fecha de inicio", nil)
	ErrDescripcionExcedeLongitud    = NewAppError("DESCRIPCION_EXCEDE_LONGITUD", "La descripción no puede exceder 1000 caracteres", nil)
)

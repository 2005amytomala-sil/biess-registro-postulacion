package domain

import (
	"fmt"
)

type AppError struct {
	Code    string
	Message string
	Err     error
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
)

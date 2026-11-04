package middleware

import (
	"errors"
	"net/http"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error  string `json:"error"`
	Codigo string `json:"codigo"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var appErr *domain.AppError
			if errors.As(err, &appErr) {
				// Handle domain errors
				status := httpStatusFromAppError(appErr)
				c.JSON(status, ErrorResponse{
					Error:  appErr.Message,
					Codigo: appErr.Code,
				})
				return
			}

			// Handle unknown errors
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:  "Error interno del servidor",
				Codigo: "INTERNAL",
			})
		}
	}
}

func httpStatusFromAppError(err *domain.AppError) int {
	switch err.Code {
	case domain.ErrNotFound.Code:
		return http.StatusNotFound
	case domain.ErrUnauthorized.Code:
		return http.StatusUnauthorized
	case domain.ErrForbidden.Code:
		return http.StatusForbidden
	case domain.ErrBadRequest.Code:
		return http.StatusBadRequest
	case domain.ErrConflict.Code:
		return http.StatusConflict
	case domain.ErrValidation.Code:
		return http.StatusBadRequest
	case domain.ErrTooManyRequests.Code:
		return http.StatusTooManyRequests
	case "YA_REGISTRADO":
		return http.StatusConflict
	case "IDENTIDAD_INVALIDA":
		return http.StatusUnauthorized
	case "CREDENCIALES_INVALIDAS":
		return http.StatusUnauthorized
	case "PASSWORD_MISMATCH":
		return http.StatusBadRequest
	case "TOKEN_INVALIDO":
		return http.StatusUnauthorized
	case "BLOQUEADO":
		return http.StatusTooManyRequests
	case "NO_AUTORIZADO":
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

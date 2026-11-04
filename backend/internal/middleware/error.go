package middleware

import (
	"errors"
	"log"
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
				if status >= 500 {
					log.Printf("[ERROR INTERNO] %s %s: %s | %v", c.Request.Method, c.Request.URL.Path, appErr.Code, appErr.Err)
				}
				body := gin.H{"error": appErr.Message, "codigo": appErr.Code}
				for k, v := range detallesDe(appErr) {
					if _, reservado := body[k]; !reservado {
						body[k] = v
					}
				}
				c.JSON(status, body)
				return
			}

			// Handle unknown errors
			log.Printf("[ERROR INTERNO] %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:  "Error interno del servidor",
				Codigo: "INTERNAL",
			})
		}
	}
}

// detallesDe busca el primer AppError de la cadena que tenga detalles
func detallesDe(err error) map[string]any {
	for err != nil {
		if appErr, ok := err.(*domain.AppError); ok && len(appErr.Detalles) > 0 {
			return appErr.Detalles
		}
		err = errors.Unwrap(err)
	}
	return nil
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
	case domain.ErrValidation.Code,
		domain.ErrFechaFormatoInvalido.Code,
		domain.ErrNacionalidadNoPermitida.Code,
		domain.ErrSectorNoValido.Code,
		domain.ErrCantonNoPerteneceAProvincia.Code,
		domain.ErrListaFamiliaresVacia.Code,
		domain.ErrFechaFinAnteriorInicio.Code,
		domain.ErrDescripcionExcedeLongitud.Code,
		domain.ErrPostulacionIncompleta.Code:
		return http.StatusUnprocessableEntity
	case domain.ErrPostulacionCerrada.Code, domain.ErrPostulacionEnviada.Code:
		return http.StatusConflict
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

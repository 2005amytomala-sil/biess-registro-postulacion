package handler

import (
	"net/http"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/middleware"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type PostulacionHandler struct {
	postulacionService *service.PostulacionService
}

func NewPostulacionHandler(postulacionService *service.PostulacionService) *PostulacionHandler {
	return &PostulacionHandler{
		postulacionService: postulacionService,
	}
}

// ObtenerPostulacion GET /api/v1/postulaciones/me
// Obtiene la postulación del usuario autenticado
func (h *PostulacionHandler) ObtenerPostulacion(c *gin.Context) {
	postulanteID, ok := middleware.PostulanteIDFrom(c)
	if !ok {
		c.Error(domain.ErrNoAutorizado)
		return
	}

	postulacion, err := h.postulacionService.ObtenerPostulacion(c.Request.Context(), postulanteID)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "No existe postulación para este usuario",
			})
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, postulacion)
}

// GuardarPostulacion PUT /api/v1/postulaciones/me
// Guarda el progreso de la postulación (secciones opcionales)
func (h *PostulacionHandler) GuardarPostulacion(c *gin.Context) {
	postulanteID, ok := middleware.PostulanteIDFrom(c)
	if !ok {
		c.Error(domain.ErrNoAutorizado)
		return
	}

	cedula, ok := middleware.CedulaFrom(c)
	if !ok {
		c.Error(domain.ErrNoAutorizado)
		return
	}

	var req domain.GuardarPostulacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(domain.ErrBadRequest)
		return
	}

	if err := h.postulacionService.GuardarPostulacion(c.Request.Context(), postulanteID, cedula, &req); err != nil {
		// Manejar errores específicos de validación de negocio
		if appErr, ok := err.(*domain.AppError); ok {
			switch appErr.Code {
			case domain.ErrNacionalidadNoPermitida.Code:
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":  appErr.Message,
					"codigo": appErr.Code,
				})
				return
			case domain.ErrSectorNoValido.Code:
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":  appErr.Message,
					"codigo": appErr.Code,
				})
				return
			case domain.ErrCantonNoPerteneceAProvincia.Code,
				domain.ErrListaFamiliaresVacia.Code,
				domain.ErrFechaFinAnteriorInicio.Code,
				domain.ErrDescripcionExcedeLongitud.Code:
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":  appErr.Message,
					"codigo": appErr.Code,
				})
				return
			}
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Postulación guardada exitosamente",
	})
}

// EnviarPostulacion POST /api/v1/postulaciones/me/enviar
// Envía la postulación (valida completitud y cambia estado a "enviada")
func (h *PostulacionHandler) EnviarPostulacion(c *gin.Context) {
	postulanteID, ok := middleware.PostulanteIDFrom(c)
	if !ok {
		c.Error(domain.ErrNoAutorizado)
		return
	}

	if err := h.postulacionService.EnviarPostulacion(c.Request.Context(), postulanteID); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Postulación enviada exitosamente",
	})
}

package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type PostulacionHandler struct {
	postulacionService domain.PostulacionService
}

func NewPostulacionHandler(postulacionService domain.PostulacionService) *PostulacionHandler {
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

	cedula, ok := middleware.CedulaFrom(c)
	if !ok {
		c.Error(domain.ErrNoAutorizado)
		return
	}

	postulacion, err := h.postulacionService.ObtenerPostulacion(c.Request.Context(), postulanteID, cedula)
	if err != nil {
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
		c.Error(errorDeBinding(err))
		return
	}

	if err := h.postulacionService.GuardarPostulacion(c.Request.Context(), postulanteID, cedula, &req); err != nil {
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

	cedula, ok := middleware.CedulaFrom(c)
	if !ok {
		c.Error(domain.ErrNoAutorizado)
		return
	}

	// El body es opcional: sin body el consentimiento se considera no aceptado
	var req domain.EnviarPostulacionRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.Error(errorDeBinding(err))
		return
	}

	codigo, err := h.postulacionService.EnviarPostulacion(c.Request.Context(), postulanteID, cedula, req.AceptaConsentimiento)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Postulación enviada exitosamente",
		"codigo":  codigo,
	})
}

// errorDeBinding conserva los errores de dominio producidos al decodificar (p. ej. formato de fecha)
func errorDeBinding(err error) error {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return domain.ErrBadRequest
}

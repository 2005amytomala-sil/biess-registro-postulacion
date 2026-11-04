package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// responderErrorAuth agrega header Retry-After y body con retry_after en bloqueos; el resto va al middleware
func responderErrorAuth(c *gin.Context, err error) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) && appErr.Code == domain.ErrBloqueado.Code {
		retryAfter := appErr.RetryAfter
		if retryAfter <= 0 {
			retryAfter = 600 // fallback por seguridad
		}
		c.Header("Retry-After", strconv.Itoa(retryAfter))
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":       appErr.Message,
			"codigo":      appErr.Code,
			"retry_after": retryAfter,
		})
		return
	}
	c.Error(err)
}

// VerificarIdentidad POST /api/v1/auth/verificar-identidad
func (h *AuthHandler) VerificarIdentidad(c *gin.Context) {
	var req domain.VerificarIdentidadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(domain.ErrBadRequest)
		return
	}

	data, err := h.authService.VerificarIdentidad(c.Request.Context(), req.Cedula, req.FechaExpedicion, req.CodigoDactilar)
	if err != nil {
		responderErrorAuth(c, err)
		return
	}

	c.JSON(http.StatusOK, data)
}

// Registrar POST /api/v1/auth/registro
func (h *AuthHandler) Registrar(c *gin.Context) {
	var req domain.RegistroRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(domain.ErrBadRequest)
		return
	}

	response, err := h.authService.Registrar(c.Request.Context(), req.Cedula, req.FechaExpedicion, req.CodigoDactilar, req.Correo, req.Password, req.ConfirmPassword)
	if err != nil {
		responderErrorAuth(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(domain.ErrBadRequest)
		return
	}

	response, err := h.authService.Login(c.Request.Context(), req.Cedula, req.Password)
	if err != nil {
		responderErrorAuth(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Recuperar POST /api/v1/auth/recuperar
func (h *AuthHandler) Recuperar(c *gin.Context) {
	var req domain.RecuperarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(domain.ErrBadRequest)
		return
	}

	if err := h.authService.Recuperar(c.Request.Context(), req.Cedula, req.FechaExpedicion, req.CodigoDactilar); err != nil {
		responderErrorAuth(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Se ha enviado un correo con las instrucciones para restablecer su contraseña",
	})
}

// Restablecer POST /api/v1/auth/restablecer
func (h *AuthHandler) Restablecer(c *gin.Context) {
	var req domain.RestablecerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(domain.ErrBadRequest)
		return
	}

	if err := h.authService.Restablecer(c.Request.Context(), req.Cedula, req.Token, req.Password, req.ConfirmPassword); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Contraseña restablecida exitosamente",
	})
}

// Me GET /api/v1/auth/me devuelve los datos del usuario autenticado
func (h *AuthHandler) Me(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"postulante_id": postulanteID,
		"cedula":        cedula,
	})
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
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
		c.Error(err)
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

	response, err := h.authService.Registrar(c.Request.Context(), req.Cedula, req.Correo, req.Password, req.ConfirmPassword)
	if err != nil {
		c.Error(err)
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
		// Si es error de bloqueo, agregar header Retry-After
		if appErr, ok := err.(*domain.AppError); ok && appErr.Code == domain.ErrBloqueado.Code {
			retryAfter := 600 // 10 minutos por defecto
			c.Header("Retry-After", strconv.Itoa(retryAfter))
		}
		c.Error(err)
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
		c.Error(err)
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

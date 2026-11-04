package handler

import (
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/middleware"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/postgres"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/redis"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	health           *HealthHandler
	auth             *AuthHandler
	catalogo         *CatalogoHandler
	postulacion      *PostulacionHandler
	authMiddleware   *middleware.AuthMiddleware
}

func NewHandler(pg *postgres.Postgres, rdb *redis.Redis, auth *AuthHandler, catalogo *CatalogoHandler, postulacion *PostulacionHandler, authMiddleware *middleware.AuthMiddleware) *Handler {
	return &Handler{
		health:         NewHealthHandler(pg, rdb),
		auth:           auth,
		catalogo:       catalogo,
		postulacion:    postulacion,
		authMiddleware: authMiddleware,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.health.Check)

	// Auth routes
	api := r.Group("/api/v1/auth")
	{
		api.POST("/verificar-identidad", h.auth.VerificarIdentidad)
		api.POST("/registro", h.auth.Registrar)
		api.POST("/login", h.auth.Login)
		api.POST("/recuperar", h.auth.Recuperar)
		api.POST("/restablecer", h.auth.Restablecer)
		api.GET("/me", h.authMiddleware.RequireAuth(), h.auth.Me)
	}

	// Catalog routes
	catalogos := r.Group("/api/v1/catalogos")
	{
		catalogos.GET("/paises", h.catalogo.ObtenerPaises)
		catalogos.GET("/provincias", h.catalogo.ObtenerProvincias)
		catalogos.GET("/provincias/:id/cantones", h.catalogo.ObtenerCantonesPorProvincia)
	}

	// Postulacion routes (protegidas con JWT)
	postulaciones := r.Group("/api/v1/postulaciones")
	postulaciones.Use(h.authMiddleware.RequireAuth())
	{
		postulaciones.GET("/me", h.postulacion.ObtenerPostulacion)
		postulaciones.PUT("/me", h.postulacion.GuardarPostulacion)
		postulaciones.POST("/me/enviar", h.postulacion.EnviarPostulacion)
	}
}

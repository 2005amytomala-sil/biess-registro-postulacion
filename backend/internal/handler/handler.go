package handler

import (
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/postgres"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/redis"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	health    *HealthHandler
	auth      *AuthHandler
	catalogo  *CatalogoHandler
}

func NewHandler(pg *postgres.Postgres, rdb *redis.Redis, auth *AuthHandler, catalogo *CatalogoHandler) *Handler {
	return &Handler{
		health:   NewHealthHandler(pg, rdb),
		auth:     auth,
		catalogo: catalogo,
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
	}

	// Catalog routes
	catalogos := r.Group("/api/v1/catalogos")
	{
		catalogos.GET("/paises", h.catalogo.ObtenerPaises)
		catalogos.GET("/provincias", h.catalogo.ObtenerProvincias)
		catalogos.GET("/provincias/:id/cantones", h.catalogo.ObtenerCantonesPorProvincia)
	}
}

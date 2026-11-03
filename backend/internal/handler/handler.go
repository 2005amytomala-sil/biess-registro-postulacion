package handler

import (
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/postgres"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/redis"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	health *HealthHandler
}

func NewHandler(pg *postgres.Postgres, redis *redis.Redis) *Handler {
	return &Handler{
		health: NewHealthHandler(pg, redis),
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.health.Check)
}

package handler

import (
	"net/http"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/postgres"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/redis"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	pg    *postgres.Postgres
	redis *redis.Redis
}

func NewHealthHandler(pg *postgres.Postgres, redis *redis.Redis) *HealthHandler {
	return &HealthHandler{
		pg:    pg,
		redis: redis,
	}
}

type HealthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

func (h *HealthHandler) Check(c *gin.Context) {
	services := make(map[string]string)
	overallStatus := "ok"

	// Check PostgreSQL
	if err := h.pg.Ping(c.Request.Context()); err != nil {
		services["postgres"] = "unhealthy"
		overallStatus = "degraded"
	} else {
		services["postgres"] = "healthy"
	}

	// Check Redis
	if err := h.redis.Ping(c.Request.Context()); err != nil {
		services["redis"] = "unhealthy"
		overallStatus = "degraded"
	} else {
		services["redis"] = "healthy"
	}

	c.JSON(http.StatusOK, HealthResponse{
		Status:   overallStatus,
		Services: services,
	})
}

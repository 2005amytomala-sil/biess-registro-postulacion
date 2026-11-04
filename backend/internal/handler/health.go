package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Pinger es cualquier dependencia que pueda verificar su conectividad (Postgres, Redis)
type Pinger interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	pg    Pinger
	redis Pinger
}

func NewHealthHandler(pg Pinger, redis Pinger) *HealthHandler {
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

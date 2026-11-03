package handler

import "github.com/gin-gonic/gin"

type Handler struct {
	// Dependencies will be injected here
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// Routes will be registered here
}

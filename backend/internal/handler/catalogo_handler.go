package handler

import (
	"net/http"
	"strconv"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type CatalogoHandler struct {
	catalogoService domain.CatalogoService
}

func NewCatalogoHandler(catalogoService domain.CatalogoService) *CatalogoHandler {
	return &CatalogoHandler{
		catalogoService: catalogoService,
	}
}

// ObtenerPaises GET /api/v1/catalogos/paises
func (h *CatalogoHandler) ObtenerPaises(c *gin.Context) {
	paises, err := h.catalogoService.ObtenerPaises(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, paises)
}

// ObtenerProvincias GET /api/v1/catalogos/provincias
func (h *CatalogoHandler) ObtenerProvincias(c *gin.Context) {
	provincias, err := h.catalogoService.ObtenerProvincias(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, provincias)
}

// ObtenerCantonesPorProvincia GET /api/v1/catalogos/provincias/:id/cantones
func (h *CatalogoHandler) ObtenerCantonesPorProvincia(c *gin.Context) {
	// Validar parámetro id
	idParam := c.Param("id")
	provinciaID, err := strconv.Atoi(idParam)
	if err != nil {
		c.Error(domain.ErrBadRequest)
		return
	}

	if provinciaID <= 0 {
		c.Error(domain.ErrBadRequest)
		return
	}

	cantones, err := h.catalogoService.ObtenerCantonesPorProvincia(c.Request.Context(), provinciaID)
	if err != nil {
		// Verificar si es error de provincia no encontrada
		if err.Error() == "provincia con ID "+strconv.Itoa(provinciaID)+" no encontrada" {
			c.Error(domain.ErrNotFound)
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cantones)
}

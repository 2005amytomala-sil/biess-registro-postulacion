package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ejecutarConError(t *testing.T, err error) (int, map[string]any) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ErrorHandler())
	r.GET("/", func(c *gin.Context) { c.Error(err) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	return w.Code, body
}

func TestErrorHandler_PostulacionCerrada409ConMotivo(t *testing.T) {
	err := domain.ErrPostulacionCerrada.ConDetalles("No es afiliado", map[string]any{"motivo_cierre": "No es afiliado", "error": "no debe pisar"})
	status, body := ejecutarConError(t, err)

	assert.Equal(t, http.StatusConflict, status)
	assert.Equal(t, "No es afiliado", body["error"])
	assert.Equal(t, "POSTULACION_CERRADA", body["codigo"])
	assert.Equal(t, "No es afiliado", body["motivo_cierre"])
}

func TestErrorHandler_Incompleta422ConFaltantes(t *testing.T) {
	err := domain.ErrPostulacionIncompleta.ConDetalles("incompleta", map[string]any{
		"faltantes": []domain.Faltante{{Seccion: "domicilio", Mensaje: "Domicilio"}},
	})
	status, body := ejecutarConError(t, fmt.Errorf("envuelto: %w", err))

	assert.Equal(t, http.StatusUnprocessableEntity, status)
	faltantes, ok := body["faltantes"].([]any)
	require.True(t, ok)
	assert.Len(t, faltantes, 1)
}

func TestErrorHandler_ReglasDeNegocio422(t *testing.T) {
	for _, err := range []error{
		domain.ErrNacionalidadNoPermitida.ConMensaje("motivo"),
		domain.ErrSectorNoValido.ConMensaje("mensaje del mock"),
		domain.ErrValidation.ConMensaje("celular"),
		domain.ErrFechaFormatoInvalido,
	} {
		status, _ := ejecutarConError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, status, err.Error())
	}
}

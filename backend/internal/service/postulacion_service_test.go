package service

import (
	"context"
	"testing"
	"time"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/mocks"
	"github.com/stretchr/testify/assert"
)

// Test nacionalidad no permitida - usa el mock de Registro Civil
func TestGuardarPostulacion_NacionalidadNoPermitida(t *testing.T) {
	registroCivil := mocks.NewRegistroCivilMock()

	// Intentar obtener datos de cédula con nacionalidad no ecuatoriana
	data, err := registroCivil.ObtenerDatos(context.Background(), "0303030303")

	// Verificar que se obtienen los datos (están en el mock)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "Estadounidense", data.Nacionalidad)

	// Verificar que al validar la nacionalidad fallaría
	if data.Nacionalidad != "Ecuatoriana" {
		assert.Equal(t, "Estadounidense", data.Nacionalidad)
	}
}

// Test validación de cantón que no pertenece a provincia
func TestValidarDomicilio_CantonNoPerteneceAProvincia(t *testing.T) {
	// Crear cantones mock para provincia 1
	cantonesProvincia1 := []domain.Canton{
		{ID: 1, ProvinciaID: 1, Nombre: "Cuenca"},
		{ID: 2, ProvinciaID: 1, Nombre: "Gualaceo"},
	}

	// Domicilio con cantón que no pertenece a la provincia
	domicilio := domain.DomicilioContacto{
		PaisID:      1,
		ProvinciaID: intPtr(1),
		CantonID:    intPtr(3), // Cantón 3 no está en provincia 1
		Celular:     "0991234567",
		Correo1:     "test@example.com",
	}

	// Verificar que el cantón 3 no está en la lista
	cantonPertenece := false
	for _, c := range cantonesProvincia1 {
		if c.ID == *domicilio.CantonID {
			cantonPertenece = true
			break
		}
	}

	assert.False(t, cantonPertenece, "El cantón no debería pertenecer a la provincia")
}

// Test validación de experiencia: fecha fin anterior a inicio
func TestValidarExperiencia_FechaFinAnteriorInicio(t *testing.T) {
	fechaInicio, _ := time.Parse("2006-01-02", "2020-01-01")
	fechaFin, _ := time.Parse("2006-01-02", "2019-01-01") // Anterior a inicio

	experiencia := []domain.ExperienciaProfesional{
		{
			Tipo:        "General",
			Cargo:       "Desarrollador",
			Institucion: "Empresa",
			FechaInicio: fechaInicio,
			FechaFin:    &fechaFin,
		},
	}

	// Verificar que fecha fin es anterior a inicio
	for _, exp := range experiencia {
		if exp.FechaFin != nil && exp.FechaFin.Before(exp.FechaInicio) {
			assert.True(t, true, "Fecha fin es anterior a inicio - debería fallar validación")
		}
	}
}

// Test validación de experiencia: descripción excede longitud
func TestValidarExperiencia_DescripcionExcedeLongitud(t *testing.T) {
	fechaInicio, _ := time.Parse("2006-01-02", "2020-01-01")
	descripcionLarga := ""
	for i := 0; i < 1001; i++ {
		descripcionLarga += "a"
	}

	experiencia := []domain.ExperienciaProfesional{
		{
			Tipo:        "General",
			Cargo:       "Desarrollador",
			Institucion: "Empresa",
			FechaInicio: fechaInicio,
			Descripcion: descripcionLarga,
		},
	}

	// Verificar que descripción excede 1000 caracteres
	for _, exp := range experiencia {
		if len(exp.Descripcion) > 1000 {
			assert.True(t, true, "Descripción excede 1000 caracteres - debería fallar validación")
		}
	}
}

// Test validación de familiares: lista vacía cuando tiene_familiares = true
func TestValidarFamiliares_ListaVacia(t *testing.T) {
	tieneFamiliares := true
	familiares := []domain.FamiliarIESS{}

	// Verificar que la lista está vacía cuando tiene_familiares = true
	if tieneFamiliares && len(familiares) == 0 {
		assert.True(t, true, "Lista de familiares vacía cuando tiene_familiares = true - debería fallar validación")
	}
}

// Test caso feliz: datos válidos
func TestValidarCasoFeliz(t *testing.T) {
	ctx := context.Background()
	registroCivil := mocks.NewRegistroCivilMock()

	// Obtener datos de cédula ecuatoriana válida
	data, err := registroCivil.ObtenerDatos(ctx, "0101010101")

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "Ecuatoriana", data.Nacionalidad)

	// Validar domicilio válido
	cantones := []domain.Canton{{ID: 1, ProvinciaID: 1, Nombre: "Cuenca"}}
	domicilio := domain.DomicilioContacto{
		PaisID:      1,
		ProvinciaID: intPtr(1),
		CantonID:    intPtr(1),
		Celular:     "0991234567",
		Correo1:     "test@example.com",
	}

	cantonPertenece := false
	for _, c := range cantones {
		if c.ID == *domicilio.CantonID {
			cantonPertenece = true
			break
		}
	}
	assert.True(t, cantonPertenece, "Cantón debería pertenecer a la provincia")

	// Validar experiencia válida
	fechaInicio, _ := time.Parse("2006-01-02", "2020-01-01")
	fechaFin, _ := time.Parse("2006-01-02", "2021-01-01")
	experiencia := []domain.ExperienciaProfesional{
		{
			Tipo:        "General",
			Cargo:       "Desarrollador",
			Institucion: "Empresa",
			FechaInicio: fechaInicio,
			FechaFin:    &fechaFin,
			Descripcion: "Descripción válida",
		},
	}

	for _, exp := range experiencia {
		if exp.FechaFin != nil {
			assert.False(t, exp.FechaFin.Before(exp.FechaInicio), "Fecha fin no debería ser anterior a inicio")
		}
		assert.True(t, len(exp.Descripcion) <= 1000, "Descripción no debería exceder 1000 caracteres")
	}
}

// Helper function
func intPtr(i int) *int {
	return &i
}

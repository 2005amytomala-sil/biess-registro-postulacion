package service

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func fecha(t *testing.T, s string) domain.Date {
	t.Helper()
	d, err := domain.ParseDate(s)
	require.NoError(t, err)
	return d
}

func assertCodigo(t *testing.T, err error, codigo string) {
	t.Helper()
	var appErr *domain.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, codigo, appErr.Code)
}

// Catálogo: Ecuador no tiene id 1 a propósito para verificar que se identifica por nombre
func catalogoMock() *MockCatalogoRepository {
	repo := &MockCatalogoRepository{}
	repo.On("ObtenerPaisPorID", mock.Anything, 7).Return(&domain.Pais{ID: 7, Nombre: "Ecuador"}, nil)
	repo.On("ObtenerPaisPorID", mock.Anything, 60).Return(&domain.Pais{ID: 60, Nombre: "España"}, nil)
	repo.On("ObtenerPaisPorID", mock.Anything, 999).Return(nil, nil)
	repo.On("ObtenerCantonesPorProvincia", mock.Anything, 19).Return([]domain.Canton{
		{ID: 100, ProvinciaID: 19, Nombre: "Quito"},
		{ID: 101, ProvinciaID: 19, Nombre: "Cayambe"},
	}, nil)
	return repo
}

func domicilioEcuador() *domain.DomicilioContacto {
	return &domain.DomicilioContacto{
		PaisID: 7, ProvinciaID: intPtr(19), CantonID: intPtr(100),
		ProvinciaTexto: "basura", Celular: "0991234567", Correo1: "test@example.com",
	}
}

func TestValidarDomicilio(t *testing.T) {
	ctx := context.Background()

	t.Run("Ecuador válido limpia campos de texto", func(t *testing.T) {
		d := domicilioEcuador()
		require.NoError(t, validarDomicilio(ctx, catalogoMock(), d))
		assert.Empty(t, d.ProvinciaTexto)
	})

	t.Run("extranjero válido limpia provincia y cantón", func(t *testing.T) {
		d := &domain.DomicilioContacto{PaisID: 60, ProvinciaID: intPtr(19), CantonID: intPtr(100),
			ProvinciaTexto: "Madrid", CiudadTexto: "Madrid", Celular: "0991234567", Correo1: "a@b.com", Correo2: "c@d.es"}
		require.NoError(t, validarDomicilio(ctx, catalogoMock(), d))
		assert.Nil(t, d.ProvinciaID)
		assert.Nil(t, d.CantonID)
	})

	casos := []struct {
		nombre string
		mutar  func(d *domain.DomicilioContacto)
		codigo string
	}{
		{"país inexistente", func(d *domain.DomicilioContacto) { d.PaisID = 999 }, "VALIDATION"},
		{"país vacío", func(d *domain.DomicilioContacto) { d.PaisID = 0 }, "VALIDATION"},
		{"Ecuador sin cantón", func(d *domain.DomicilioContacto) { d.CantonID = nil }, "VALIDATION"},
		{"cantón de otra provincia", func(d *domain.DomicilioContacto) { d.CantonID = intPtr(5) }, domain.ErrCantonNoPerteneceAProvincia.Code},
		{"extranjero sin ciudad", func(d *domain.DomicilioContacto) { d.PaisID = 60; d.ProvinciaTexto = "Madrid"; d.CiudadTexto = " " }, "VALIDATION"},
		{"celular con 9 dígitos", func(d *domain.DomicilioContacto) { d.Celular = "099123456" }, "VALIDATION"},
		{"celular con letras", func(d *domain.DomicilioContacto) { d.Celular = "09912345ab" }, "VALIDATION"},
		{"correo 1 inválido", func(d *domain.DomicilioContacto) { d.Correo1 = "no-es-correo" }, "VALIDATION"},
		{"correo 1 vacío", func(d *domain.DomicilioContacto) { d.Correo1 = "" }, "VALIDATION"},
		{"correo 2 inválido", func(d *domain.DomicilioContacto) { d.Correo2 = "x@y" }, "VALIDATION"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			d := domicilioEcuador()
			c.mutar(d)
			assertCodigo(t, validarDomicilio(ctx, catalogoMock(), d), c.codigo)
		})
	}
}

func TestValidarRespaldo(t *testing.T) {
	assertCodigo(t, validarRespaldo(&domain.RespaldoOrganizacion{TieneRespaldo: true, NombreOrganizacion: "  "}), "VALIDATION")

	r := &domain.RespaldoOrganizacion{TieneRespaldo: true, NombreOrganizacion: " Gremio "}
	require.NoError(t, validarRespaldo(r))
	assert.Equal(t, "Gremio", r.NombreOrganizacion)

	r = &domain.RespaldoOrganizacion{TieneRespaldo: false, NombreOrganizacion: "sobra"}
	require.NoError(t, validarRespaldo(r))
	assert.Empty(t, r.NombreOrganizacion)
}

func TestValidarFamiliares(t *testing.T) {
	assert.Equal(t, domain.ErrListaFamiliaresVacia, validarFamiliares(true, nil))
	require.NoError(t, validarFamiliares(false, nil))
	assertCodigo(t, validarFamiliares(true, []domain.FamiliarIESS{{Nombres: "Ana", Parentesco: "Hermana", Institucion: "SRI"}}), "VALIDATION")
	require.NoError(t, validarFamiliares(true, []domain.FamiliarIESS{{Nombres: "Ana", Parentesco: "Hermana", Institucion: "BIESS"}}))
}

func TestValidarExperiencia(t *testing.T) {
	base := func() domain.ExperienciaProfesional {
		return domain.ExperienciaProfesional{Tipo: "General", Cargo: "Analista", Institucion: "IESS", FechaInicio: fecha(t, "2020-01-01")}
	}

	exp := base()
	require.NoError(t, validarExperiencia([]domain.ExperienciaProfesional{exp}), "sin fecha fin es válido")

	exp.FechaFin = fecha(t, "2020-01-01")
	require.NoError(t, validarExperiencia([]domain.ExperienciaProfesional{exp}), "fin igual a inicio es válido")

	exp.FechaFin = fecha(t, "2019-12-31")
	assert.Equal(t, domain.ErrFechaFinAnteriorInicio, validarExperiencia([]domain.ExperienciaProfesional{exp}))

	exp = base()
	exp.Descripcion = strings.Repeat("ñ", 1000)
	require.NoError(t, validarExperiencia([]domain.ExperienciaProfesional{exp}), "1000 caracteres multibyte es válido")
	exp.Descripcion += "a"
	assert.Equal(t, domain.ErrDescripcionExcedeLongitud, validarExperiencia([]domain.ExperienciaProfesional{exp}))

	exp = base()
	exp.FechaInicio = domain.Date{}
	assertCodigo(t, validarExperiencia([]domain.ExperienciaProfesional{exp}), "VALIDATION")

	exp = base()
	exp.Tipo = "Otro"
	assertCodigo(t, validarExperiencia([]domain.ExperienciaProfesional{exp}), "VALIDATION")
}

func TestValidarFormacion(t *testing.T) {
	ok := domain.FormacionAcademica{Nivel: "Tercer nivel", Institucion: "UCE", Tipo: "Nacional", Titulo: "Ing."}
	require.NoError(t, validarFormacion([]domain.FormacionAcademica{ok}))

	ok.Tipo = "Otra"
	assertCodigo(t, validarFormacion([]domain.FormacionAcademica{ok}), "VALIDATION")
}

func TestValidarInhabilidades(t *testing.T) {
	assertCodigo(t, validarInhabilidades([]domain.Inhabilidad{{CodigoPregunta: "1.9.14"}}), "VALIDATION")
	assertCodigo(t, validarInhabilidades([]domain.Inhabilidad{{CodigoPregunta: "1.9.1"}, {CodigoPregunta: "1.9.1"}}), "VALIDATION")
	assertCodigo(t, validarInhabilidades([]domain.Inhabilidad{{CodigoPregunta: "1.9.7", Respuesta: true}}), "VALIDATION")
	assertCodigo(t, validarInhabilidades([]domain.Inhabilidad{{CodigoPregunta: "1.9.7", Respuesta: true, Detalle: `{"institucion":"IESS","cargo":"","fecha_finalizacion":"2024-01-01","motivo":"x"}`}}), "VALIDATION")

	inh := []domain.Inhabilidad{{CodigoPregunta: "1.9.7", Respuesta: false}}
	require.NoError(t, validarInhabilidades(inh))
	assert.Equal(t, "{}", inh[0].Detalle, "1.9.7 en No debe guardar detalle no nulo (CHECK de la BD)")

	require.NoError(t, validarInhabilidades([]domain.Inhabilidad{{CodigoPregunta: "1.9.7", Respuesta: true,
		Detalle: `{"institucion":"IESS","cargo":"Director","fecha_finalizacion":"2024-01-01","motivo":"Renuncia"}`}}))
}

func TestValidarRequest_SoloSeccionesPresentes(t *testing.T) {
	// Sin catálogo configurado: si validara domicilio ausente, el mock entraría en pánico
	repo := &MockCatalogoRepository{}
	req := &domain.GuardarPostulacionRequest{
		Inhabilidades: []domain.Inhabilidad{{CodigoPregunta: "1.9.1"}},
	}
	require.NoError(t, validarRequest(context.Background(), repo, req))

	req = &domain.GuardarPostulacionRequest{TieneFamiliares: boolPtr(true)}
	assert.Equal(t, domain.ErrListaFamiliaresVacia, validarRequest(context.Background(), repo, req))
}

func TestSectorEfectivo(t *testing.T) {
	borrador := &domain.Postulacion{Sector: "jubilados", Estado: domain.EstadoBorrador}
	cerrada := &domain.Postulacion{Sector: "jubilados", Estado: domain.EstadoCerrada}

	assert.Equal(t, "afiliados", sectorEfectivo("", nil))
	assert.Equal(t, "empleadores", sectorEfectivo("empleadores", nil))
	assert.Equal(t, "jubilados", sectorEfectivo("", borrador))
	assert.Equal(t, "empleadores", sectorEfectivo("empleadores", borrador))
	assert.Equal(t, "jubilados", sectorEfectivo("empleadores", cerrada), "una postulación no editable no cambia de sector")
}

func TestVerificarEditable(t *testing.T) {
	require.NoError(t, verificarEditable(&domain.Postulacion{Estado: domain.EstadoBorrador}))

	err := verificarEditable(&domain.Postulacion{Estado: domain.EstadoCerrada, MotivoCierre: "Motivo X"})
	var appErr *domain.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domain.ErrPostulacionCerrada.Code, appErr.Code)
	assert.Equal(t, "Motivo X", appErr.Message)
	assert.Equal(t, "Motivo X", appErr.Detalles["motivo_cierre"])

	assertCodigo(t, verificarEditable(&domain.Postulacion{Estado: domain.EstadoEnviada}), domain.ErrPostulacionEnviada.Code)
}

func TestEvaluarCompletitud(t *testing.T) {
	completo := &domain.ResumenCompletitud{
		DatosPersonales: true, Domicilio: true, Respaldo: true,
		InhabilidadesRespondidas: append([]string(nil), domain.CodigosInhabilidades...),
	}
	assert.Empty(t, evaluarCompletitud(completo, "afiliados", true))

	secciones := func(f []domain.Faltante) []string {
		var s []string
		for _, x := range f {
			s = append(s, x.Seccion)
		}
		return s
	}

	assert.Equal(t,
		[]string{"datos_personales", "domicilio", "respaldo_organizacion", "inhabilidades", "consentimiento"},
		secciones(evaluarCompletitud(&domain.ResumenCompletitud{}, "afiliados", false)))

	parcial := *completo
	parcial.InhabilidadesRespondidas = []string{"1.9.1", "1.9.2"}
	parcial.TieneRespaldo = true
	faltantes := evaluarCompletitud(&parcial, "afiliados", true)
	require.Len(t, faltantes, 2)
	assert.Equal(t, "respaldo_organizacion", faltantes[0].Seccion)
	assert.Contains(t, faltantes[1].Mensaje, "1.9.3")
	assert.NotContains(t, faltantes[1].Mensaje, "1.9.2,")
}

func TestDate_JSON(t *testing.T) {
	var exp domain.ExperienciaProfesional
	require.NoError(t, json.Unmarshal([]byte(`{"fecha_inicio":"2020-05-01","fecha_fin":""}`), &exp))
	assert.True(t, exp.FechaInicio.Valid)
	assert.Equal(t, "2020-05-01", exp.FechaInicio.String())
	assert.False(t, exp.FechaFin.Valid)

	require.NoError(t, json.Unmarshal([]byte(`{"fecha_inicio":"2020-05-01","fecha_fin":null}`), &exp))
	assert.False(t, exp.FechaFin.Valid)

	err := json.Unmarshal([]byte(`{"fecha_inicio":"01/05/2020"}`), &exp)
	assertCodigo(t, err, domain.ErrFechaFormatoInvalido.Code)

	out, err := json.Marshal(domain.ExperienciaProfesional{FechaInicio: fecha(t, "2021-02-03")})
	require.NoError(t, err)
	assert.Contains(t, string(out), `"fecha_inicio":"2021-02-03"`)
	assert.Contains(t, string(out), `"fecha_fin":null`)
}

func TestDate_ValueScan(t *testing.T) {
	v, err := domain.Date{}.Value()
	require.NoError(t, err)
	assert.Nil(t, v)

	var d domain.Date
	require.NoError(t, d.Scan(nil))
	assert.False(t, d.Valid)
	require.NoError(t, d.Scan(fecha(t, "2022-07-08").Time))
	assert.Equal(t, "2022-07-08", d.String())
}

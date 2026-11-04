package service

import (
	"context"
	"errors"
	"testing"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// fakeTxManager simula BEGIN/COMMIT/ROLLBACK: cuenta cuántas transacciones terminaron en cada estado
type fakeTxManager struct {
	commits   int
	rollbacks int
}

func (f *fakeTxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if err := fn(ctx); err != nil {
		f.rollbacks++
		return err
	}
	f.commits++
	return nil
}

type mockPostulacionRepo struct{ mock.Mock }

func (m *mockPostulacionRepo) ObtenerPorPostulanteID(ctx context.Context, id int) (*domain.Postulacion, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*domain.Postulacion)
	return p, args.Error(1)
}
func (m *mockPostulacionRepo) ObtenerPorPostulanteIDParaActualizar(ctx context.Context, id int) (*domain.Postulacion, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*domain.Postulacion)
	return p, args.Error(1)
}
func (m *mockPostulacionRepo) Crear(ctx context.Context, p *domain.Postulacion) error {
	return m.Called(ctx, p).Error(0)
}
func (m *mockPostulacionRepo) Actualizar(ctx context.Context, p *domain.Postulacion) error {
	return m.Called(ctx, p).Error(0)
}
func (m *mockPostulacionRepo) GenerarCodigoUnico(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}
func (m *mockPostulacionRepo) ResumenCompletitud(ctx context.Context, id int) (*domain.ResumenCompletitud, error) {
	args := m.Called(ctx, id)
	r, _ := args.Get(0).(*domain.ResumenCompletitud)
	return r, args.Error(1)
}

// mockSeccion implementa todos los repositorios de secciones; cada método registra la llamada
type mockSeccion struct{ mock.Mock }

func (m *mockSeccion) EliminarPorPostulacionID(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockSeccion) guardar(ctx context.Context, v any) error { return m.Called(ctx, v).Error(0) }

type mockDatosPersonales struct{ mockSeccion }

func (m *mockDatosPersonales) Guardar(ctx context.Context, d *domain.DatosPersonales) error {
	return m.guardar(ctx, d)
}
func (m *mockDatosPersonales) ObtenerPorPostulacionID(context.Context, int) (*domain.DatosPersonales, error) {
	return nil, nil
}

type mockDomicilio struct{ mockSeccion }

func (m *mockDomicilio) Guardar(ctx context.Context, d *domain.DomicilioContacto) error {
	return m.guardar(ctx, d)
}
func (m *mockDomicilio) ObtenerPorPostulacionID(context.Context, int) (*domain.DomicilioContacto, error) {
	return nil, nil
}

type mockRespaldo struct{ mockSeccion }

func (m *mockRespaldo) Guardar(ctx context.Context, r *domain.RespaldoOrganizacion) error {
	return m.guardar(ctx, r)
}
func (m *mockRespaldo) ObtenerPorPostulacionID(context.Context, int) (*domain.RespaldoOrganizacion, error) {
	return nil, nil
}

type mockFamiliares struct{ mockSeccion }

func (m *mockFamiliares) Guardar(ctx context.Context, f []domain.FamiliarIESS) error {
	return m.guardar(ctx, f)
}
func (m *mockFamiliares) ObtenerPorPostulacionID(context.Context, int) ([]domain.FamiliarIESS, error) {
	return nil, nil
}

type mockFormacion struct{ mockSeccion }

func (m *mockFormacion) Guardar(ctx context.Context, f []domain.FormacionAcademica) error {
	return m.guardar(ctx, f)
}
func (m *mockFormacion) ObtenerPorPostulacionID(context.Context, int) ([]domain.FormacionAcademica, error) {
	return nil, nil
}

type mockExperiencia struct{ mockSeccion }

func (m *mockExperiencia) Guardar(ctx context.Context, e []domain.ExperienciaProfesional) error {
	return m.guardar(ctx, e)
}
func (m *mockExperiencia) ObtenerPorPostulacionID(context.Context, int) ([]domain.ExperienciaProfesional, error) {
	return nil, nil
}

type mockInhabilidades struct{ mockSeccion }

func (m *mockInhabilidades) Guardar(ctx context.Context, i []domain.Inhabilidad) error {
	return m.guardar(ctx, i)
}
func (m *mockInhabilidades) ObtenerPorPostulacionID(context.Context, int) ([]domain.Inhabilidad, error) {
	return nil, nil
}

type escenario struct {
	tx            *fakeTxManager
	postulaciones *mockPostulacionRepo
	datos         *mockDatosPersonales
	domicilio     *mockDomicilio
	familiares    *mockFamiliares
	respaldo      *mockRespaldo
	formacion     *mockFormacion
	experiencia   *mockExperiencia
	inhabilidades *mockInhabilidades
	catalogo      *MockCatalogoRepository
	service       *PostulacionService
}

const postulanteID = 42

func nuevoEscenario() *escenario {
	e := &escenario{
		tx:            &fakeTxManager{},
		postulaciones: &mockPostulacionRepo{},
		datos:         &mockDatosPersonales{},
		domicilio:     &mockDomicilio{},
		familiares:    &mockFamiliares{},
		respaldo:      &mockRespaldo{},
		formacion:     &mockFormacion{},
		experiencia:   &mockExperiencia{},
		inhabilidades: &mockInhabilidades{},
		catalogo:      catalogoMock(),
	}
	e.service = NewPostulacionService(e.tx, e.postulaciones, e.datos, e.domicilio, e.familiares, e.respaldo,
		e.formacion, e.experiencia, e.inhabilidades, mocks.NewRegistroCivilMock(), mocks.NewHistorialLaboralMock(), e.catalogo)
	return e
}

// conBorrador simula una postulación existente en borrador
func (e *escenario) conBorrador() *domain.Postulacion {
	p := &domain.Postulacion{ID: 10, PostulanteID: postulanteID, Codigo: "POST-1", Sector: "afiliados", Estado: domain.EstadoBorrador}
	e.postulaciones.On("ObtenerPorPostulanteID", mock.Anything, postulanteID).Return(p, nil)
	e.postulaciones.On("ObtenerPorPostulanteIDParaActualizar", mock.Anything, postulanteID).Return(p, nil)
	return p
}

func TestGuardarPostulacion_NacionalidadNoEcuatoriana_CierrePersistido(t *testing.T) {
	e := nuevoEscenario()
	e.conBorrador()
	e.postulaciones.On("Actualizar", mock.Anything, mock.MatchedBy(func(p *domain.Postulacion) bool {
		return p.Estado == domain.EstadoCerrada && p.MotivoCierre == motivoNacionalidad
	})).Return(nil).Once()

	req := &domain.GuardarPostulacionRequest{Sector: "afiliados", Domicilio: domicilioEcuador()}
	err := e.service.GuardarPostulacion(context.Background(), postulanteID, "0303030303", req)

	require.ErrorIs(t, err, domain.ErrNacionalidadNoPermitida)
	assertCodigo(t, err, domain.ErrNacionalidadNoPermitida.Code)
	assert.Equal(t, 1, e.tx.commits, "el cierre debe confirmarse en su propia transacción")
	assert.Equal(t, 0, e.tx.rollbacks)
	e.postulaciones.AssertExpectations(t)
	e.domicilio.AssertNotCalled(t, "guardar", mock.Anything, mock.Anything)
	e.datos.AssertNotCalled(t, "guardar", mock.Anything, mock.Anything)
}

func TestGuardarPostulacion_NacionalidadNoEcuatoriana_SinPostulacionPreviaLaCreaCerrada(t *testing.T) {
	e := nuevoEscenario()
	e.postulaciones.On("ObtenerPorPostulanteID", mock.Anything, postulanteID).Return(nil, domain.ErrNotFound)
	e.postulaciones.On("ObtenerPorPostulanteIDParaActualizar", mock.Anything, postulanteID).Return(nil, domain.ErrNotFound)
	e.postulaciones.On("GenerarCodigoUnico", mock.Anything).Return("POST-X", nil)
	e.postulaciones.On("Crear", mock.Anything, mock.MatchedBy(func(p *domain.Postulacion) bool {
		return p.Estado == domain.EstadoCerrada && p.MotivoCierre == motivoNacionalidad
	})).Return(nil).Once()

	_, err := e.service.ObtenerPostulacion(context.Background(), postulanteID, "0303030303")

	require.ErrorIs(t, err, domain.ErrNacionalidadNoPermitida)
	assert.Equal(t, 1, e.tx.commits)
	e.postulaciones.AssertExpectations(t)
}

func TestGuardarPostulacion_SectorNoValido_CierreConMensajeDelMock(t *testing.T) {
	e := nuevoEscenario()
	e.conBorrador()
	e.postulaciones.On("Actualizar", mock.Anything, mock.MatchedBy(func(p *domain.Postulacion) bool {
		return p.Estado == domain.EstadoCerrada && p.Sector == "jubilados" && p.MotivoCierre != ""
	})).Return(nil).Once()

	err := e.service.GuardarPostulacion(context.Background(), postulanteID, "0404040404", &domain.GuardarPostulacionRequest{Sector: "jubilados"})

	var appErr *domain.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domain.ErrSectorNoValido.Code, appErr.Code)
	_, mensaje, _ := mocks.NewHistorialLaboralMock().VerificarSector(context.Background(), "0404040404", "jubilados")
	assert.Equal(t, mensaje, appErr.Message)
	assert.Equal(t, 1, e.tx.commits)
	e.postulaciones.AssertExpectations(t)
}

func TestGuardarPostulacion_CantonNoPerteneceAProvincia(t *testing.T) {
	e := nuevoEscenario()
	e.conBorrador()

	domicilio := domicilioEcuador()
	domicilio.CantonID = intPtr(5) // no pertenece a la provincia 19
	err := e.service.GuardarPostulacion(context.Background(), postulanteID, "0101010101",
		&domain.GuardarPostulacionRequest{Sector: "afiliados", Domicilio: domicilio})

	require.ErrorIs(t, err, domain.ErrCantonNoPerteneceAProvincia)
	var appErr *domain.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, "canton_id", appErr.Detalles["campo"])
	assert.Equal(t, 0, e.tx.commits+e.tx.rollbacks, "no debe abrirse la transacción de guardado")
	e.domicilio.AssertNotCalled(t, "guardar", mock.Anything, mock.Anything)
}

func TestGuardarPostulacion_RollbackSiFallaUnaInsercion(t *testing.T) {
	e := nuevoEscenario()
	e.conBorrador()
	errBD := errors.New("violación de restricción")
	e.datos.On("guardar", mock.Anything, mock.Anything).Return(nil)
	e.domicilio.On("guardar", mock.Anything, mock.Anything).Return(nil)
	e.experiencia.On("EliminarPorPostulacionID", mock.Anything, 10).Return(nil)
	e.experiencia.On("guardar", mock.Anything, mock.Anything).Return(errBD)

	req := &domain.GuardarPostulacionRequest{
		Sector:    "afiliados",
		Domicilio: domicilioEcuador(),
		Experiencia: []domain.ExperienciaProfesional{
			{Tipo: "General", Cargo: "Analista", Institucion: "IESS", FechaInicio: fecha(t, "2020-01-01")},
		},
		Inhabilidades: []domain.Inhabilidad{{CodigoPregunta: "1.9.1"}},
	}
	err := e.service.GuardarPostulacion(context.Background(), postulanteID, "0101010101", req)

	require.ErrorIs(t, err, errBD)
	assert.Equal(t, 1, e.tx.rollbacks, "la transacción debe hacer ROLLBACK")
	assert.Equal(t, 0, e.tx.commits)
	e.domicilio.AssertCalled(t, "guardar", mock.Anything, mock.Anything)
	e.inhabilidades.AssertNotCalled(t, "guardar", mock.Anything, mock.Anything)
}

func TestGuardarPostulacion_GuardadoParcialSoloTocaSusTablas(t *testing.T) {
	e := nuevoEscenario()
	e.conBorrador()
	e.datos.On("guardar", mock.Anything, mock.Anything).Return(nil)
	e.inhabilidades.On("guardar", mock.Anything, mock.Anything).Return(nil)

	err := e.service.GuardarPostulacion(context.Background(), postulanteID, "0101010101",
		&domain.GuardarPostulacionRequest{Inhabilidades: []domain.Inhabilidad{{CodigoPregunta: "1.9.1"}}})

	require.NoError(t, err)
	assert.Equal(t, 1, e.tx.commits)
	e.inhabilidades.AssertExpectations(t)
	e.domicilio.AssertNotCalled(t, "guardar", mock.Anything, mock.Anything)
	e.familiares.AssertNotCalled(t, "EliminarPorPostulacionID", mock.Anything, mock.Anything)
	e.formacion.AssertNotCalled(t, "EliminarPorPostulacionID", mock.Anything, mock.Anything)
	e.experiencia.AssertNotCalled(t, "EliminarPorPostulacionID", mock.Anything, mock.Anything)
}

func TestGuardarPostulacion_Cerrada409ConMotivo(t *testing.T) {
	e := nuevoEscenario()
	cerrada := &domain.Postulacion{ID: 10, Sector: "afiliados", Estado: domain.EstadoCerrada, MotivoCierre: "Motivo original"}
	e.postulaciones.On("ObtenerPorPostulanteID", mock.Anything, postulanteID).Return(cerrada, nil)

	err := e.service.GuardarPostulacion(context.Background(), postulanteID, "0101010101", &domain.GuardarPostulacionRequest{})

	var appErr *domain.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domain.ErrPostulacionCerrada.Code, appErr.Code)
	assert.Equal(t, "Motivo original", appErr.Detalles["motivo_cierre"])
	assert.Equal(t, 0, e.tx.commits+e.tx.rollbacks)
}

func TestEnviarPostulacion_Incompleta(t *testing.T) {
	e := nuevoEscenario()
	e.conBorrador()
	e.datos.On("guardar", mock.Anything, mock.Anything).Return(nil)
	e.postulaciones.On("ResumenCompletitud", mock.Anything, 10).Return(&domain.ResumenCompletitud{
		DatosPersonales: true, Respaldo: true, InhabilidadesRespondidas: []string{"1.9.1"},
	}, nil)

	codigo, err := e.service.EnviarPostulacion(context.Background(), postulanteID, "0101010101", false)

	assert.Empty(t, codigo)
	var appErr *domain.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domain.ErrPostulacionIncompleta.Code, appErr.Code)
	faltantes, ok := appErr.Detalles["faltantes"].([]domain.Faltante)
	require.True(t, ok)
	var secciones []string
	for _, f := range faltantes {
		secciones = append(secciones, f.Seccion)
	}
	assert.Equal(t, []string{"domicilio", "inhabilidades", "consentimiento"}, secciones)
	assert.Equal(t, 1, e.tx.rollbacks)
	e.postulaciones.AssertNotCalled(t, "Actualizar", mock.Anything, mock.Anything)
}

func TestEnviarPostulacion_Completa(t *testing.T) {
	e := nuevoEscenario()
	e.conBorrador()
	e.datos.On("guardar", mock.Anything, mock.Anything).Return(nil)
	e.postulaciones.On("ResumenCompletitud", mock.Anything, 10).Return(&domain.ResumenCompletitud{
		DatosPersonales: true, Domicilio: true, Respaldo: true,
		InhabilidadesRespondidas: append([]string(nil), domain.CodigosInhabilidades...),
	}, nil)
	e.postulaciones.On("GenerarCodigoUnico", mock.Anything).Return("POST-FINAL", nil)
	e.postulaciones.On("Actualizar", mock.Anything, mock.MatchedBy(func(p *domain.Postulacion) bool {
		return p.Estado == domain.EstadoEnviada && p.Codigo == "POST-FINAL"
	})).Return(nil).Once()

	codigo, err := e.service.EnviarPostulacion(context.Background(), postulanteID, "0101010101", true)

	require.NoError(t, err)
	assert.Equal(t, "POST-FINAL", codigo)
	assert.Equal(t, 1, e.tx.commits)
	e.postulaciones.AssertExpectations(t)
}

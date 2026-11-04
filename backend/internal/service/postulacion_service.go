package service

import (
	"context"
	"fmt"
	"time"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/postgres"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/repository"
	"github.com/jackc/pgx/v5"
)

type PostulacionService struct {
	pg                     *postgres.Postgres
	postulacionRepo        *repository.PostulacionRepository
	datosPersonalesRepo    *repository.DatosPersonalesRepository
	domicilioRepo          *repository.DomicilioRepository
	familiaresRepo         *repository.FamiliaresRepository
	respaldosRepo          *repository.RespaldoRepository
	formacionRepo          *repository.FormacionRepository
	experienciaRepo        *repository.ExperienciaRepository
	inhabilidadesRepo      *repository.InhabilidadesRepository
	registroCivil          domain.RegistroCivilClient
	historialLaboral       domain.HistorialLaboralClient
	catalogoRepo           domain.CatalogoRepository
}

func NewPostulacionService(
	pg *postgres.Postgres,
	postulacionRepo *repository.PostulacionRepository,
	datosPersonalesRepo *repository.DatosPersonalesRepository,
	domicilioRepo *repository.DomicilioRepository,
	familiaresRepo *repository.FamiliaresRepository,
	respaldosRepo *repository.RespaldoRepository,
	formacionRepo *repository.FormacionRepository,
	experienciaRepo *repository.ExperienciaRepository,
	inhabilidadesRepo *repository.InhabilidadesRepository,
	registroCivil domain.RegistroCivilClient,
	historialLaboral domain.HistorialLaboralClient,
	catalogoRepo domain.CatalogoRepository,
) *PostulacionService {
	return &PostulacionService{
		pg:                     pg,
		postulacionRepo:        postulacionRepo,
		datosPersonalesRepo:    datosPersonalesRepo,
		domicilioRepo:          domicilioRepo,
		familiaresRepo:         familiaresRepo,
		respaldosRepo:          respaldosRepo,
		formacionRepo:          formacionRepo,
		experienciaRepo:        experienciaRepo,
		inhabilidadesRepo:      inhabilidadesRepo,
		registroCivil:          registroCivil,
		historialLaboral:       historialLaboral,
		catalogoRepo:           catalogoRepo,
	}
}

// GuardarPostulacion guarda o actualiza una postulación dentro de una transacción
func (s *PostulacionService) GuardarPostulacion(ctx context.Context, postulanteID int, cedula string, req *domain.GuardarPostulacionRequest) error {
	// Iniciar transacción
	if err := s.pg.WithTx(ctx, func(ctx context.Context, tx interface{}) error {
		pgxTx, ok := tx.(pgx.Tx)
		if !ok {
			return fmt.Errorf("error: transaction is not pgx.Tx")
		}

		// 1. Obtener o crear postulación
		postulacion, err := s.postulacionRepo.ObtenerPorPostulanteID(ctx, postulanteID)
		if err != nil {
			// Si no existe, crear nueva
			codigo, err := s.postulacionRepo.GenerarCodigoUnico(ctx)
			if err != nil {
				return fmt.Errorf("error generando código único: %w", err)
			}

			postulacion = &domain.Postulacion{
				PostulanteID: postulanteID,
				Codigo:        codigo,
				Sector:        req.Sector,
				Estado:        "borrador",
			}

			if err := s.postulacionRepo.Crear(ctx, pgxTx, postulacion); err != nil {
				return fmt.Errorf("error creando postulación: %w", err)
			}
		} else {
			// Actualizar sector si cambió
			if postulacion.Sector != req.Sector {
				postulacion.Sector = req.Sector
				if err := s.postulacionRepo.Actualizar(ctx, pgxTx, postulacion); err != nil {
					return fmt.Errorf("error actualizando postulación: %w", err)
				}
			}
		}

		// 2. Obtener datos personales desde Registro Civil (usando cédula del JWT)
		registroCivilData, err := s.registroCivil.ObtenerDatos(ctx, cedula)
		if err != nil {
			return fmt.Errorf("error obteniendo datos del registro civil: %w", err)
		}

		// 3. Validar nacionalidad
		if registroCivilData.Nacionalidad != "Ecuatoriana" {
			// Cerrar postulación
			postulacion.Estado = "cerrada"
			postulacion.MotivoCierre = "Nacionalidad no permitida. Solo se permite nacionalidad ecuatoriana."
			if err := s.postulacionRepo.Actualizar(ctx, pgxTx, postulacion); err != nil {
				return fmt.Errorf("error cerrando postulación: %w", err)
			}
			return domain.ErrNacionalidadNoPermitida
		}

		// 4. Validar sector con Historial Laboral
		valido, mensaje, err := s.historialLaboral.VerificarSector(ctx, cedula, req.Sector)
		if err != nil {
			return fmt.Errorf("error verificando sector: %w", err)
		}
		if !valido {
			// Cerrar postulación
			postulacion.Estado = "cerrada"
			postulacion.MotivoCierre = mensaje
			if err := s.postulacionRepo.Actualizar(ctx, pgxTx, postulacion); err != nil {
				return fmt.Errorf("error cerrando postulación: %w", err)
			}
			return &domain.AppError{
				Code:    domain.ErrSectorNoValido.Code,
				Message: mensaje,
				Err:     domain.ErrSectorNoValido,
			}
		}

		// 5. Validar domicilio geográfico
		if err := s.validarDomicilio(ctx, req.Domicilio); err != nil {
			return err
		}

		// 6. Validar familiares
		if err := s.validarFamiliares(req.TieneFamiliares, req.Familiares); err != nil {
			return err
		}

		// 7. Validar experiencia
		if err := s.validarExperiencia(req.Experiencia); err != nil {
			return err
		}

		// 8. Guardar datos personales (desde Registro Civil)
		fechaNacimiento, err := time.Parse("02/01/2006", registroCivilData.FechaNacimiento)
		if err != nil {
			return fmt.Errorf("error parseando fecha de nacimiento: %w", err)
		}

		datosPersonales := &domain.DatosPersonales{
			PostulacionID:   postulacion.ID,
			Nombres:         registroCivilData.Nombres,
			Apellidos:       registroCivilData.Apellidos,
			Cedula:          registroCivilData.Cedula,
			Genero:          registroCivilData.Genero,
			FechaNacimiento: fechaNacimiento,
			Nacionalidad:    registroCivilData.Nacionalidad,
			EstadoCivil:     registroCivilData.EstadoCivil,
			ParejaNombres:   registroCivilData.ParejaNombres,
			ParejaCedula:    registroCivilData.ParejaCedula,
		}

		if err := s.datosPersonalesRepo.Guardar(ctx, pgxTx, datosPersonales); err != nil {
			return fmt.Errorf("error guardando datos personales: %w", err)
		}

		// 9. Guardar domicilio
		domicilio := req.Domicilio
		domicilio.PostulacionID = postulacion.ID
		if err := s.domicilioRepo.Guardar(ctx, pgxTx, &domicilio); err != nil {
			return fmt.Errorf("error guardando domicilio: %w", err)
		}

		// 10. Guardar respaldo de organización
		respaldo := req.RespaldoOrganizacion
		respaldo.PostulacionID = postulacion.ID
		if err := s.respaldosRepo.Guardar(ctx, pgxTx, &respaldo); err != nil {
			return fmt.Errorf("error guardando respaldo: %w", err)
		}

		// 11. Guardar familiares (delete + insert)
		if err := s.familiaresRepo.EliminarPorPostulacionID(ctx, pgxTx, postulacion.ID); err != nil {
			return fmt.Errorf("error eliminando familiares: %w", err)
		}

		if len(req.Familiares) > 0 {
			for i := range req.Familiares {
				req.Familiares[i].PostulacionID = postulacion.ID
			}
			if err := s.familiaresRepo.Guardar(ctx, pgxTx, req.Familiares); err != nil {
				return fmt.Errorf("error guardando familiares: %w", err)
			}
		}

		// 12. Guardar formación (delete + insert)
		if err := s.formacionRepo.EliminarPorPostulacionID(ctx, pgxTx, postulacion.ID); err != nil {
			return fmt.Errorf("error eliminando formación: %w", err)
		}

		if len(req.Formacion) > 0 {
			for i := range req.Formacion {
				req.Formacion[i].PostulacionID = postulacion.ID
			}
			if err := s.formacionRepo.Guardar(ctx, pgxTx, req.Formacion); err != nil {
				return fmt.Errorf("error guardando formación: %w", err)
			}
		}

		// 13. Guardar experiencia (delete + insert)
		if err := s.experienciaRepo.EliminarPorPostulacionID(ctx, pgxTx, postulacion.ID); err != nil {
			return fmt.Errorf("error eliminando experiencia: %w", err)
		}

		if len(req.Experiencia) > 0 {
			for i := range req.Experiencia {
				req.Experiencia[i].PostulacionID = postulacion.ID
			}
			if err := s.experienciaRepo.Guardar(ctx, pgxTx, req.Experiencia); err != nil {
				return fmt.Errorf("error guardando experiencia: %w", err)
			}
		}

		// 14. Guardar inhabilidades (upsert)
		if len(req.Inhabilidades) > 0 {
			for i := range req.Inhabilidades {
				req.Inhabilidades[i].PostulacionID = postulacion.ID
			}
			if err := s.inhabilidadesRepo.Guardar(ctx, pgxTx, req.Inhabilidades); err != nil {
				return fmt.Errorf("error guardando inhabilidades: %w", err)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// validarDomicilio valida las reglas geográficas del domicilio
func (s *PostulacionService) validarDomicilio(ctx context.Context, domicilio domain.DomicilioContacto) error {
	// Si es Ecuador (pais_id = 1), provincia y cantón son obligatorios
	if domicilio.PaisID == 1 {
		if domicilio.ProvinciaID == nil || domicilio.CantonID == nil {
			return domain.ErrValidation
		}

		// Validar que el cantón pertenece a la provincia
		cantones, err := s.catalogoRepo.ObtenerCantonesPorProvincia(ctx, *domicilio.ProvinciaID)
		if err != nil {
			return fmt.Errorf("error obteniendo cantones: %w", err)
		}

		cantonPertenece := false
		for _, c := range cantones {
			if c.ID == *domicilio.CantonID {
				cantonPertenece = true
				break
			}
		}

		if !cantonPertenece {
			return domain.ErrCantonNoPerteneceAProvincia
		}
	} else {
		// Si no es Ecuador, los campos de texto deben estar presentes
		if domicilio.ProvinciaTexto == "" || domicilio.CiudadTexto == "" {
			return domain.ErrValidation
		}
	}

	return nil
}

// validarFamiliares valida la lista de familiares
func (s *PostulacionService) validarFamiliares(tieneFamiliares bool, familiares []domain.FamiliarIESS) error {
	if tieneFamiliares {
		if len(familiares) == 0 {
			return domain.ErrListaFamiliaresVacia
		}

		// Validar cada familiar
		for _, familiar := range familiares {
			if familiar.Nombres == "" || familiar.Parentesco == "" || familiar.Institucion == "" {
				return domain.ErrValidation
			}
		}
	}

	return nil
}

// validarExperiencia valida la experiencia profesional
func (s *PostulacionService) validarExperiencia(experiencia []domain.ExperienciaProfesional) error {
	for _, exp := range experiencia {
		// Validar fecha fin >= fecha inicio
		if exp.FechaFin != nil && exp.FechaFin.Before(exp.FechaInicio) {
			return domain.ErrFechaFinAnteriorInicio
		}

		// Validar descripción <= 1000 caracteres
		if len(exp.Descripcion) > 1000 {
			return domain.ErrDescripcionExcedeLongitud
		}
	}

	return nil
}

// ObtenerPostulacion obtiene la postulación de un postulante
func (s *PostulacionService) ObtenerPostulacion(ctx context.Context, postulanteID int) (*domain.Postulacion, error) {
	postulacion, err := s.postulacionRepo.ObtenerPorPostulanteID(ctx, postulanteID)
	if err != nil {
		return nil, err
	}

	return postulacion, nil
}

// EnviarPostulacion envía la postulación (cambia estado a "enviada" y genera código único)
func (s *PostulacionService) EnviarPostulacion(ctx context.Context, postulanteID int) error {
	if err := s.pg.WithTx(ctx, func(ctx context.Context, tx interface{}) error {
		pgxTx, ok := tx.(pgx.Tx)
		if !ok {
			return fmt.Errorf("error: transaction is not pgx.Tx")
		}

		// Obtener postulación
		postulacion, obtainErr := s.postulacionRepo.ObtenerPorPostulanteID(ctx, postulanteID)
		if obtainErr != nil {
			return fmt.Errorf("error obteniendo postulación: %w", obtainErr)
		}

		// Validar que esté en estado borrador
		if postulacion.Estado != "borrador" {
			return &domain.AppError{
				Code:    "ESTADO_INVALIDO",
				Message: "Solo se puede enviar una postulación en estado borrador",
				Err:     domain.ErrValidation,
			}
		}

		// Generar código único definitivo
		codigo, err := s.postulacionRepo.GenerarCodigoUnico(ctx)
		if err != nil {
			return fmt.Errorf("error generando código único: %w", err)
		}

		// Actualizar estado y código
		postulacion.Estado = "enviada"
		postulacion.Codigo = codigo
		if err := s.postulacionRepo.Actualizar(ctx, pgxTx, postulacion); err != nil {
			return fmt.Errorf("error actualizando postulación: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

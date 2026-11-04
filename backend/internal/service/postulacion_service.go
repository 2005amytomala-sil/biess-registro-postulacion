package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
)

// PostulacionService implementa domain.PostulacionService; depende solo de interfaces de domain
type PostulacionService struct {
	tx                  domain.TxManager
	postulacionRepo     domain.PostulacionRepository
	datosPersonalesRepo domain.DatosPersonalesRepository
	domicilioRepo       domain.DomicilioRepository
	familiaresRepo      domain.FamiliaresRepository
	respaldosRepo       domain.RespaldoRepository
	formacionRepo       domain.FormacionRepository
	experienciaRepo     domain.ExperienciaRepository
	inhabilidadesRepo   domain.InhabilidadesRepository
	registroCivil       domain.RegistroCivilClient
	historialLaboral    domain.HistorialLaboralClient
	catalogoRepo        domain.CatalogoRepository
}

func NewPostulacionService(
	tx domain.TxManager,
	postulacionRepo domain.PostulacionRepository,
	datosPersonalesRepo domain.DatosPersonalesRepository,
	domicilioRepo domain.DomicilioRepository,
	familiaresRepo domain.FamiliaresRepository,
	respaldosRepo domain.RespaldoRepository,
	formacionRepo domain.FormacionRepository,
	experienciaRepo domain.ExperienciaRepository,
	inhabilidadesRepo domain.InhabilidadesRepository,
	registroCivil domain.RegistroCivilClient,
	historialLaboral domain.HistorialLaboralClient,
	catalogoRepo domain.CatalogoRepository,
) *PostulacionService {
	return &PostulacionService{
		tx:                  tx,
		postulacionRepo:     postulacionRepo,
		datosPersonalesRepo: datosPersonalesRepo,
		domicilioRepo:       domicilioRepo,
		familiaresRepo:      familiaresRepo,
		respaldosRepo:       respaldosRepo,
		formacionRepo:       formacionRepo,
		experienciaRepo:     experienciaRepo,
		inhabilidadesRepo:   inhabilidadesRepo,
		registroCivil:       registroCivil,
		historialLaboral:    historialLaboral,
		catalogoRepo:        catalogoRepo,
	}
}

// elegibilidad agrupa el resultado de las verificaciones externas previas a cualquier operación
type elegibilidad struct {
	registroCivil *domain.RegistroCivilData
	sector        string
	existente     *domain.Postulacion
}

func (s *PostulacionService) obtenerExistente(ctx context.Context, postulanteID int) (*domain.Postulacion, error) {
	p, err := s.postulacionRepo.ObtenerPorPostulanteID(ctx, postulanteID)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, nil
	}
	return p, err
}

// verificarElegibilidad consulta el Registro Civil (cédula del JWT) y el Historial Laboral.
// Si la nacionalidad o el sector no son válidos, cierra la postulación en una transacción propia
// (con COMMIT) y recién después devuelve el error.
func (s *PostulacionService) verificarElegibilidad(ctx context.Context, postulanteID int, cedula string, sectorSolicitado string) (*elegibilidad, error) {
	rc, err := s.registroCivil.ObtenerDatos(ctx, cedula)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos del registro civil: %w", err)
	}

	existente, err := s.obtenerExistente(ctx, postulanteID)
	if err != nil {
		return nil, err
	}

	if rc.Nacionalidad != nacionalidadPermitida {
		if err := s.cerrarPostulacion(ctx, postulanteID, "", motivoNacionalidad); err != nil {
			return nil, err
		}
		return nil, domain.ErrNacionalidadNoPermitida.ConMensaje(motivoNacionalidad)
	}

	sector := sectorEfectivo(sectorSolicitado, existente)
	if !sectorValido(sector) {
		return nil, errValidacion("El sector debe ser afiliados, jubilados o empleadores")
	}

	valido, mensaje, err := s.historialLaboral.VerificarSector(ctx, cedula, sector)
	if err != nil {
		return nil, fmt.Errorf("error verificando sector: %w", err)
	}
	if !valido {
		if err := s.cerrarPostulacion(ctx, postulanteID, sector, mensaje); err != nil {
			return nil, err
		}
		return nil, domain.ErrSectorNoValido.ConMensaje(mensaje)
	}

	return &elegibilidad{registroCivil: rc, sector: sector, existente: existente}, nil
}

// cerrarPostulacion persiste estado=cerrada y motivo_cierre en su propia transacción.
// Solo cierra postulaciones en borrador (no altera una enviada ni el motivo de una ya cerrada).
func (s *PostulacionService) cerrarPostulacion(ctx context.Context, postulanteID int, sector string, motivo string) error {
	return s.tx.WithTx(ctx, func(ctx context.Context) error {
		postulacion, err := s.postulacionRepo.ObtenerPorPostulanteIDParaActualizar(ctx, postulanteID)
		if errors.Is(err, domain.ErrNotFound) {
			codigo, err := s.postulacionRepo.GenerarCodigoUnico(ctx)
			if err != nil {
				return fmt.Errorf("error generando código único: %w", err)
			}
			if sector == "" {
				sector = sectorPorDefecto
			}
			return s.postulacionRepo.Crear(ctx, &domain.Postulacion{
				PostulanteID: postulanteID,
				Codigo:       codigo,
				Sector:       sector,
				Estado:       domain.EstadoCerrada,
				MotivoCierre: motivo,
			})
		}
		if err != nil {
			return err
		}
		if postulacion.Estado != domain.EstadoBorrador {
			return nil
		}

		postulacion.Estado = domain.EstadoCerrada
		postulacion.MotivoCierre = motivo
		if sector != "" {
			postulacion.Sector = sector
		}
		if err := s.postulacionRepo.Actualizar(ctx, postulacion); err != nil {
			return fmt.Errorf("error cerrando postulación: %w", err)
		}
		return nil
	})
}

// obtenerOCrear obtiene la postulación bloqueando la fila, o la crea en borrador
func (s *PostulacionService) obtenerOCrear(ctx context.Context, postulanteID int, sector string) (*domain.Postulacion, error) {
	postulacion, err := s.postulacionRepo.ObtenerPorPostulanteIDParaActualizar(ctx, postulanteID)
	if err == nil {
		return postulacion, nil
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	codigo, err := s.postulacionRepo.GenerarCodigoUnico(ctx)
	if err != nil {
		return nil, fmt.Errorf("error generando código único: %w", err)
	}
	postulacion = &domain.Postulacion{
		PostulanteID: postulanteID,
		Codigo:       codigo,
		Sector:       sector,
		Estado:       domain.EstadoBorrador,
	}
	if err := s.postulacionRepo.Crear(ctx, postulacion); err != nil {
		return nil, fmt.Errorf("error creando postulación: %w", err)
	}
	return postulacion, nil
}

// guardarDatosPersonales persiste los datos personales tomados del Registro Civil
func (s *PostulacionService) guardarDatosPersonales(ctx context.Context, postulacionID int, rc *domain.RegistroCivilData) error {
	fechaNacimiento, err := time.Parse("02/01/2006", rc.FechaNacimiento)
	if err != nil {
		return fmt.Errorf("error parseando fecha de nacimiento: %w", err)
	}

	datosPersonales := &domain.DatosPersonales{
		PostulacionID:   postulacionID,
		Nombres:         rc.Nombres,
		Apellidos:       rc.Apellidos,
		Cedula:          rc.Cedula,
		Genero:          rc.Genero,
		FechaNacimiento: fechaNacimiento,
		Nacionalidad:    rc.Nacionalidad,
		EstadoCivil:     rc.EstadoCivil,
		ParejaNombres:   rc.ParejaNombres,
		ParejaCedula:    rc.ParejaCedula,
	}
	if err := s.datosPersonalesRepo.Guardar(ctx, datosPersonales); err != nil {
		return fmt.Errorf("error guardando datos personales: %w", err)
	}
	return nil
}

// GuardarPostulacion guarda de forma parcial: cada sección presente en el request modifica solo su tabla
func (s *PostulacionService) GuardarPostulacion(ctx context.Context, postulanteID int, cedula string, req *domain.GuardarPostulacionRequest) error {
	eleg, err := s.verificarElegibilidad(ctx, postulanteID, cedula, req.Sector)
	if err != nil {
		return err
	}
	if eleg.existente != nil {
		if err := verificarEditable(eleg.existente); err != nil {
			return err
		}
	}
	if err := validarRequest(ctx, s.catalogoRepo, req); err != nil {
		return err
	}

	return s.tx.WithTx(ctx, func(ctx context.Context) error {
		postulacion, err := s.obtenerOCrear(ctx, postulanteID, eleg.sector)
		if err != nil {
			return err
		}
		// Revalidar con la fila bloqueada por si cambió el estado de forma concurrente
		if err := verificarEditable(postulacion); err != nil {
			return err
		}
		if postulacion.Sector != eleg.sector {
			postulacion.Sector = eleg.sector
			if err := s.postulacionRepo.Actualizar(ctx, postulacion); err != nil {
				return fmt.Errorf("error actualizando postulación: %w", err)
			}
		}

		if err := s.guardarDatosPersonales(ctx, postulacion.ID, eleg.registroCivil); err != nil {
			return err
		}

		if req.Domicilio != nil {
			req.Domicilio.PostulacionID = postulacion.ID
			if err := s.domicilioRepo.Guardar(ctx, req.Domicilio); err != nil {
				return fmt.Errorf("error guardando domicilio: %w", err)
			}
		}

		if req.TieneFamiliares != nil || req.Familiares != nil {
			if err := s.familiaresRepo.EliminarPorPostulacionID(ctx, postulacion.ID); err != nil {
				return fmt.Errorf("error eliminando familiares: %w", err)
			}
			if tieneFamiliares(req) {
				for i := range req.Familiares {
					req.Familiares[i].PostulacionID = postulacion.ID
				}
				if err := s.familiaresRepo.Guardar(ctx, req.Familiares); err != nil {
					return fmt.Errorf("error guardando familiares: %w", err)
				}
			}
		}

		if req.RespaldoOrganizacion != nil {
			req.RespaldoOrganizacion.PostulacionID = postulacion.ID
			if err := s.respaldosRepo.Guardar(ctx, req.RespaldoOrganizacion); err != nil {
				return fmt.Errorf("error guardando respaldo: %w", err)
			}
		}

		if req.Formacion != nil {
			if err := s.formacionRepo.EliminarPorPostulacionID(ctx, postulacion.ID); err != nil {
				return fmt.Errorf("error eliminando formación: %w", err)
			}
			for i := range req.Formacion {
				req.Formacion[i].PostulacionID = postulacion.ID
			}
			if err := s.formacionRepo.Guardar(ctx, req.Formacion); err != nil {
				return fmt.Errorf("error guardando formación: %w", err)
			}
		}

		if req.Experiencia != nil {
			if err := s.experienciaRepo.EliminarPorPostulacionID(ctx, postulacion.ID); err != nil {
				return fmt.Errorf("error eliminando experiencia: %w", err)
			}
			for i := range req.Experiencia {
				req.Experiencia[i].PostulacionID = postulacion.ID
			}
			if err := s.experienciaRepo.Guardar(ctx, req.Experiencia); err != nil {
				return fmt.Errorf("error guardando experiencia: %w", err)
			}
		}

		if req.Inhabilidades != nil {
			for i := range req.Inhabilidades {
				req.Inhabilidades[i].PostulacionID = postulacion.ID
			}
			if err := s.inhabilidadesRepo.Guardar(ctx, req.Inhabilidades); err != nil {
				return fmt.Errorf("error guardando inhabilidades: %w", err)
			}
		}

		return nil
	})
}

// ObtenerPostulacion obtiene la postulación de un postulante, creándola si no existe
func (s *PostulacionService) ObtenerPostulacion(ctx context.Context, postulanteID int, cedula string) (map[string]interface{}, error) {
	eleg, err := s.verificarElegibilidad(ctx, postulanteID, cedula, "")
	if err != nil {
		return nil, err
	}

	postulacion := eleg.existente
	if postulacion == nil {
		codigo, err := s.postulacionRepo.GenerarCodigoUnico(ctx)
		if err != nil {
			return nil, fmt.Errorf("error generando código único: %w", err)
		}

		postulacion = &domain.Postulacion{
			PostulanteID: postulanteID,
			Codigo:       codigo,
			Sector:       eleg.sector,
			Estado:       domain.EstadoBorrador,
		}

		if err := s.postulacionRepo.Crear(ctx, postulacion); err != nil {
			return nil, fmt.Errorf("error creando postulación: %w", err)
		}
	}
	registroCivilData := eleg.registroCivil
	// Obtener domicilio si existe
	domicilio, err := s.domicilioRepo.ObtenerPorPostulacionID(ctx, postulacion.ID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo domicilio: %w", err)
	}

	// Obtener familiares si existen
	familiares, err := s.familiaresRepo.ObtenerPorPostulacionID(ctx, postulacion.ID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo familiares: %w", err)
	}

	// Obtener respaldo si existe
	respaldo, err := s.respaldosRepo.ObtenerPorPostulacionID(ctx, postulacion.ID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo respaldo: %w", err)
	}

	// Obtener formación si existe
	formacion, err := s.formacionRepo.ObtenerPorPostulacionID(ctx, postulacion.ID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo formación: %w", err)
	}

	// Obtener experiencia si existe
	experiencia, err := s.experienciaRepo.ObtenerPorPostulacionID(ctx, postulacion.ID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo experiencia: %w", err)
	}

	// Obtener inhabilidades si existen
	inhabilidades, err := s.inhabilidadesRepo.ObtenerPorPostulacionID(ctx, postulacion.ID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo inhabilidades: %w", err)
	}

	// Construir respuesta completa
	response := map[string]interface{}{
		"id":                    postulacion.ID,
		"postulante_id":         postulacion.PostulanteID,
		"codigo":                postulacion.Codigo,
		"sector":                postulacion.Sector,
		"estado":                postulacion.Estado,
		"motivo_cierre":         postulacion.MotivoCierre,
		"created_at":            postulacion.CreatedAt,
		"updated_at":            postulacion.UpdatedAt,
		"datos_personales":      registroCivilData,
		"domicilio":             domicilio,
		"tiene_familiares":      len(familiares) > 0,
		"familiares":            familiares,
		"respaldo_organizacion": respaldo,
		"formacion":             formacion,
		"experiencia":           experiencia,
		"inhabilidades":         inhabilidades,
	}

	return response, nil
}

// EnviarPostulacion verifica elegibilidad y completitud y cambia el estado a "enviada", todo en una transacción
func (s *PostulacionService) EnviarPostulacion(ctx context.Context, postulanteID int, cedula string, aceptaConsentimiento bool) (string, error) {
	eleg, err := s.verificarElegibilidad(ctx, postulanteID, cedula, "")
	if err != nil {
		return "", err
	}

	var codigo string
	err = s.tx.WithTx(ctx, func(ctx context.Context) error {
		postulacion, err := s.postulacionRepo.ObtenerPorPostulanteIDParaActualizar(ctx, postulanteID)
		if errors.Is(err, domain.ErrNotFound) {
			faltantes := evaluarCompletitud(&domain.ResumenCompletitud{}, eleg.sector, aceptaConsentimiento)
			return domain.ErrPostulacionIncompleta.ConDetalles(domain.ErrPostulacionIncompleta.Message, map[string]any{"faltantes": faltantes})
		}
		if err != nil {
			return fmt.Errorf("error obteniendo postulación: %w", err)
		}
		if err := verificarEditable(postulacion); err != nil {
			return err
		}

		// Los datos personales siempre se toman del Registro Civil vigente
		if err := s.guardarDatosPersonales(ctx, postulacion.ID, eleg.registroCivil); err != nil {
			return err
		}

		resumen, err := s.postulacionRepo.ResumenCompletitud(ctx, postulacion.ID)
		if err != nil {
			return err
		}
		if faltantes := evaluarCompletitud(resumen, postulacion.Sector, aceptaConsentimiento); len(faltantes) > 0 {
			return domain.ErrPostulacionIncompleta.ConDetalles(domain.ErrPostulacionIncompleta.Message, map[string]any{"faltantes": faltantes})
		}

		codigoGenerado, err := s.postulacionRepo.GenerarCodigoUnico(ctx)
		if err != nil {
			return fmt.Errorf("error generando código único: %w", err)
		}

		postulacion.Estado = domain.EstadoEnviada
		postulacion.Codigo = codigoGenerado
		if err := s.postulacionRepo.Actualizar(ctx, postulacion); err != nil {
			return fmt.Errorf("error actualizando postulación: %w", err)
		}
		codigo = codigoGenerado
		return nil
	})
	if err != nil {
		return "", err
	}

	return codigo, nil
}

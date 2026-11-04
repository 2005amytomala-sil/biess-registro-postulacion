package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
)

const (
	nacionalidadPermitida = "Ecuatoriana"
	motivoNacionalidad    = "Nacionalidad no permitida. Solo se permite nacionalidad ecuatoriana."
	sectorPorDefecto      = "afiliados"
	maxDescripcion        = 1000
)

var (
	regexCelular = regexp.MustCompile(`^\d{10}$`)
	regexCorreo  = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

	tiposFormacion     = []string{"Nacional", "Extranjera"}
	tiposExperiencia   = []string{"General", "Direccion"}
	institucionesIESS  = []string{"IESS", "BIESS"}
	camposDetalle1_9_7 = []string{"institucion", "cargo", "fecha_finalizacion", "motivo"}
)

func errValidacion(format string, args ...any) error {
	return domain.ErrValidation.ConMensaje(fmt.Sprintf(format, args...))
}

// errCampo es un error de validación que indica qué campo del formulario lo causó
func errCampo(campo, mensaje string) error {
	return domain.ErrValidation.ConDetalles(mensaje, map[string]any{"campo": campo})
}

func sectorValido(sector string) bool {
	return slices.Contains(domain.SectoresValidos, sector)
}

func esEcuador(pais *domain.Pais) bool {
	return pais != nil && strings.EqualFold(strings.TrimSpace(pais.Nombre), "Ecuador")
}

// sectorEfectivo decide qué sector validar: el de una postulación no editable no puede cambiarse
func sectorEfectivo(solicitado string, existente *domain.Postulacion) string {
	switch {
	case existente != nil && existente.Estado != domain.EstadoBorrador:
		return existente.Sector
	case solicitado != "":
		return solicitado
	case existente != nil:
		return existente.Sector
	default:
		return sectorPorDefecto
	}
}

// verificarEditable devuelve 409 si la postulación está cerrada o enviada
func verificarEditable(p *domain.Postulacion) error {
	detalles := map[string]any{"estado": p.Estado, "motivo_cierre": p.MotivoCierre}
	switch p.Estado {
	case domain.EstadoCerrada:
		mensaje := p.MotivoCierre
		if mensaje == "" {
			mensaje = domain.ErrPostulacionCerrada.Message
		}
		return domain.ErrPostulacionCerrada.ConDetalles(mensaje, detalles)
	case domain.EstadoEnviada:
		detalles["codigo"] = p.Codigo
		return domain.ErrPostulacionEnviada.ConDetalles(domain.ErrPostulacionEnviada.Message, detalles)
	}
	return nil
}

// validarDomicilio aplica las reglas geográficas y de contacto, y normaliza los campos según el país
func validarDomicilio(ctx context.Context, catalogo domain.CatalogoRepository, d *domain.DomicilioContacto) error {
	if d.PaisID <= 0 {
		return errCampo("pais_id", "El país es obligatorio")
	}
	pais, err := catalogo.ObtenerPaisPorID(ctx, d.PaisID)
	if err != nil {
		return err
	}
	if pais == nil {
		return errCampo("pais_id", "El país seleccionado no existe en el catálogo")
	}

	if esEcuador(pais) {
		if d.ProvinciaID == nil {
			return errCampo("provincia_id", "La provincia es obligatoria para domicilios en Ecuador")
		}
		if d.CantonID == nil {
			return errCampo("canton_id", "El cantón es obligatorio para domicilios en Ecuador")
		}
		cantones, err := catalogo.ObtenerCantonesPorProvincia(ctx, *d.ProvinciaID)
		if err != nil {
			return fmt.Errorf("error obteniendo cantones: %w", err)
		}
		if !slices.ContainsFunc(cantones, func(c domain.Canton) bool { return c.ID == *d.CantonID }) {
			return domain.ErrCantonNoPerteneceAProvincia.ConDetalles(domain.ErrCantonNoPerteneceAProvincia.Message, map[string]any{"campo": "canton_id"})
		}
		d.ProvinciaTexto, d.CiudadTexto = "", ""
	} else {
		d.ProvinciaTexto, d.CiudadTexto = strings.TrimSpace(d.ProvinciaTexto), strings.TrimSpace(d.CiudadTexto)
		if d.ProvinciaTexto == "" {
			return errCampo("provincia_texto", "La provincia/estado es obligatoria para domicilios fuera de Ecuador")
		}
		if d.CiudadTexto == "" {
			return errCampo("ciudad_texto", "La ciudad es obligatoria para domicilios fuera de Ecuador")
		}
		d.ProvinciaID, d.CantonID = nil, nil
	}

	if !regexCelular.MatchString(d.Celular) {
		return errCampo("celular", "El celular debe tener exactamente 10 dígitos")
	}
	d.Correo1, d.Correo2 = strings.TrimSpace(d.Correo1), strings.TrimSpace(d.Correo2)
	if !regexCorreo.MatchString(d.Correo1) {
		return errCampo("correo_1", "El correo electrónico principal no es válido")
	}
	if d.Correo2 != "" && !regexCorreo.MatchString(d.Correo2) {
		return errCampo("correo_2", "El correo electrónico secundario no es válido")
	}
	return nil
}

func validarFamiliares(tieneFamiliares bool, familiares []domain.FamiliarIESS) error {
	if !tieneFamiliares {
		return nil
	}
	if len(familiares) == 0 {
		return domain.ErrListaFamiliaresVacia
	}
	for _, f := range familiares {
		if strings.TrimSpace(f.Nombres) == "" || strings.TrimSpace(f.Parentesco) == "" {
			return errValidacion("Nombres y parentesco del familiar son obligatorios")
		}
		if !slices.Contains(institucionesIESS, f.Institucion) {
			return errValidacion("La institución del familiar debe ser IESS o BIESS")
		}
	}
	return nil
}

func validarRespaldo(r *domain.RespaldoOrganizacion) error {
	r.NombreOrganizacion = strings.TrimSpace(r.NombreOrganizacion)
	if r.TieneRespaldo && r.NombreOrganizacion == "" {
		return errValidacion("El nombre de la organización es obligatorio cuando cuenta con respaldo")
	}
	if !r.TieneRespaldo {
		r.NombreOrganizacion = ""
	}
	return nil
}

func validarFormacion(formacion []domain.FormacionAcademica) error {
	for _, f := range formacion {
		if strings.TrimSpace(f.Nivel) == "" || strings.TrimSpace(f.Institucion) == "" || strings.TrimSpace(f.Titulo) == "" {
			return errValidacion("Nivel, institución y título de la formación son obligatorios")
		}
		if !slices.Contains(tiposFormacion, f.Tipo) {
			return errValidacion("El tipo de formación debe ser Nacional o Extranjera")
		}
	}
	return nil
}

func validarExperiencia(experiencia []domain.ExperienciaProfesional) error {
	for _, exp := range experiencia {
		if !slices.Contains(tiposExperiencia, exp.Tipo) {
			return errValidacion("El tipo de experiencia debe ser General o Direccion")
		}
		if strings.TrimSpace(exp.Cargo) == "" || strings.TrimSpace(exp.Institucion) == "" {
			return errValidacion("Cargo e institución de la experiencia son obligatorios")
		}
		if !exp.FechaInicio.Valid {
			return errValidacion("La fecha de inicio de la experiencia es obligatoria")
		}
		if exp.FechaFin.Before(exp.FechaInicio) {
			return domain.ErrFechaFinAnteriorInicio
		}
		if utf8.RuneCountInString(exp.Descripcion) > maxDescripcion {
			return domain.ErrDescripcionExcedeLongitud
		}
	}
	return nil
}

// validarInhabilidades valida códigos y el detalle obligatorio de la 1.9.7; normaliza el detalle vacío
func validarInhabilidades(inhabilidades []domain.Inhabilidad) error {
	vistos := make(map[string]bool, len(inhabilidades))
	for i := range inhabilidades {
		inh := &inhabilidades[i]
		if !slices.Contains(domain.CodigosInhabilidades, inh.CodigoPregunta) {
			return errValidacion("Código de inhabilidad no válido: %s", inh.CodigoPregunta)
		}
		if vistos[inh.CodigoPregunta] {
			return errValidacion("La inhabilidad %s está repetida", inh.CodigoPregunta)
		}
		vistos[inh.CodigoPregunta] = true

		if inh.Detalle != "" && !json.Valid([]byte(inh.Detalle)) {
			return errValidacion("El detalle de la inhabilidad %s no es válido", inh.CodigoPregunta)
		}
		if inh.CodigoPregunta != domain.CodigoInhabilidadConDetalle {
			continue
		}
		if !inh.Respuesta {
			if inh.Detalle == "" {
				inh.Detalle = "{}"
			}
			continue
		}
		var detalle map[string]any
		if err := json.Unmarshal([]byte(inh.Detalle), &detalle); err != nil {
			return errValidacion("La inhabilidad %s requiere detalle", inh.CodigoPregunta)
		}
		for _, campo := range camposDetalle1_9_7 {
			if v, _ := detalle[campo].(string); strings.TrimSpace(v) == "" {
				return errValidacion("La inhabilidad %s requiere el campo %s", inh.CodigoPregunta, campo)
			}
		}
	}
	return nil
}

// validarRequest valida todas las secciones presentes antes de abrir la transacción
func validarRequest(ctx context.Context, catalogo domain.CatalogoRepository, req *domain.GuardarPostulacionRequest) error {
	if req.Domicilio != nil {
		if err := validarDomicilio(ctx, catalogo, req.Domicilio); err != nil {
			return err
		}
	}
	if req.TieneFamiliares != nil || req.Familiares != nil {
		if err := validarFamiliares(tieneFamiliares(req), req.Familiares); err != nil {
			return err
		}
	}
	if req.RespaldoOrganizacion != nil {
		if err := validarRespaldo(req.RespaldoOrganizacion); err != nil {
			return err
		}
	}
	if err := validarFormacion(req.Formacion); err != nil {
		return err
	}
	if err := validarExperiencia(req.Experiencia); err != nil {
		return err
	}
	return validarInhabilidades(req.Inhabilidades)
}

func tieneFamiliares(req *domain.GuardarPostulacionRequest) bool {
	if req.TieneFamiliares != nil {
		return *req.TieneFamiliares
	}
	return len(req.Familiares) > 0
}

// evaluarCompletitud lista lo que falta para poder enviar la postulación
func evaluarCompletitud(r *domain.ResumenCompletitud, sector string, aceptaConsentimiento bool) []domain.Faltante {
	var faltantes []domain.Faltante
	agregar := func(seccion, mensaje string) {
		faltantes = append(faltantes, domain.Faltante{Seccion: seccion, Mensaje: mensaje})
	}

	if !r.DatosPersonales {
		agregar("datos_personales", "Datos personales del Registro Civil")
	}
	if !r.Domicilio {
		agregar("domicilio", "Domicilio y datos de contacto")
	}
	if !sectorValido(sector) {
		agregar("sector", "Sector al que postula")
	}
	switch {
	case !r.Respaldo:
		agregar("respaldo_organizacion", "Respaldo de organización")
	case r.TieneRespaldo && strings.TrimSpace(r.NombreOrganizacion) == "":
		agregar("respaldo_organizacion", "Nombre de la organización que respalda la postulación")
	}

	var sinResponder []string
	for _, codigo := range domain.CodigosInhabilidades {
		if !slices.Contains(r.InhabilidadesRespondidas, codigo) {
			sinResponder = append(sinResponder, codigo)
		}
	}
	if len(sinResponder) > 0 {
		agregar("inhabilidades", "Inhabilidades sin responder: "+strings.Join(sinResponder, ", "))
	}
	if !aceptaConsentimiento {
		agregar("consentimiento", "Aceptación del consentimiento y declaración")
	}
	return faltantes
}

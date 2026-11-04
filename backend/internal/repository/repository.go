package repository

import "github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"

// Verificación en compilación de que cada repositorio implementa su interfaz de dominio
var (
	_ domain.PostulanteRepository      = (*PostulanteRepository)(nil)
	_ domain.CatalogoRepository        = (*CatalogoRepository)(nil)
	_ domain.PostulacionRepository     = (*PostulacionRepository)(nil)
	_ domain.DatosPersonalesRepository = (*DatosPersonalesRepository)(nil)
	_ domain.DomicilioRepository       = (*DomicilioRepository)(nil)
	_ domain.FamiliaresRepository      = (*FamiliaresRepository)(nil)
	_ domain.RespaldoRepository        = (*RespaldoRepository)(nil)
	_ domain.FormacionRepository       = (*FormacionRepository)(nil)
	_ domain.ExperienciaRepository     = (*ExperienciaRepository)(nil)
	_ domain.InhabilidadesRepository   = (*InhabilidadesRepository)(nil)
)

package domain

import "context"

// TxManager ejecuta fn dentro de una transacción. La transacción viaja en el contexto que recibe fn:
// los repositorios que usen ese contexto participan de ella. Si fn devuelve error se hace ROLLBACK.
type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// PostulacionRepository: los métodos de lectura devuelven un error que envuelve ErrNotFound si no existe
type PostulacionRepository interface {
	ObtenerPorPostulanteID(ctx context.Context, postulanteID int) (*Postulacion, error)
	// ObtenerPorPostulanteIDParaActualizar bloquea la fila (FOR UPDATE); usar dentro de WithTx
	ObtenerPorPostulanteIDParaActualizar(ctx context.Context, postulanteID int) (*Postulacion, error)
	Crear(ctx context.Context, postulacion *Postulacion) error
	Actualizar(ctx context.Context, postulacion *Postulacion) error
	GenerarCodigoUnico(ctx context.Context) (string, error)
	ResumenCompletitud(ctx context.Context, postulacionID int) (*ResumenCompletitud, error)
}

// Los repositorios de secciones de una sola fila devuelven (nil, nil) si la sección no existe

type DatosPersonalesRepository interface {
	Guardar(ctx context.Context, datos *DatosPersonales) error
	ObtenerPorPostulacionID(ctx context.Context, postulacionID int) (*DatosPersonales, error)
}

type DomicilioRepository interface {
	Guardar(ctx context.Context, domicilio *DomicilioContacto) error
	ObtenerPorPostulacionID(ctx context.Context, postulacionID int) (*DomicilioContacto, error)
}

type RespaldoRepository interface {
	Guardar(ctx context.Context, respaldo *RespaldoOrganizacion) error
	ObtenerPorPostulacionID(ctx context.Context, postulacionID int) (*RespaldoOrganizacion, error)
}

type FamiliaresRepository interface {
	EliminarPorPostulacionID(ctx context.Context, postulacionID int) error
	Guardar(ctx context.Context, familiares []FamiliarIESS) error
	ObtenerPorPostulacionID(ctx context.Context, postulacionID int) ([]FamiliarIESS, error)
}

type FormacionRepository interface {
	EliminarPorPostulacionID(ctx context.Context, postulacionID int) error
	Guardar(ctx context.Context, formacion []FormacionAcademica) error
	ObtenerPorPostulacionID(ctx context.Context, postulacionID int) ([]FormacionAcademica, error)
}

type ExperienciaRepository interface {
	EliminarPorPostulacionID(ctx context.Context, postulacionID int) error
	Guardar(ctx context.Context, experiencia []ExperienciaProfesional) error
	ObtenerPorPostulacionID(ctx context.Context, postulacionID int) ([]ExperienciaProfesional, error)
}

type InhabilidadesRepository interface {
	Guardar(ctx context.Context, inhabilidades []Inhabilidad) error
	ObtenerPorPostulacionID(ctx context.Context, postulacionID int) ([]Inhabilidad, error)
}

// Casos de uso que consumen los handlers

type AuthService interface {
	VerificarIdentidad(ctx context.Context, cedula, fechaExpedicion, codigoDactilar string) (*RegistroCivilData, error)
	Registrar(ctx context.Context, cedula, fechaExpedicion, codigoDactilar, correo, password, confirmPassword string) (*AuthResponse, error)
	Login(ctx context.Context, cedula, password string) (*AuthResponse, error)
	Recuperar(ctx context.Context, cedula, fechaExpedicion, codigoDactilar string) error
	Restablecer(ctx context.Context, cedula, token, password, confirmPassword string) error
}

type CatalogoService interface {
	ObtenerPaises(ctx context.Context) ([]Pais, error)
	ObtenerProvincias(ctx context.Context) ([]Provincia, error)
	ObtenerCantonesPorProvincia(ctx context.Context, provinciaID int) ([]Canton, error)
}

type PostulacionService interface {
	ObtenerPostulacion(ctx context.Context, postulanteID int, cedula string) (map[string]interface{}, error)
	GuardarPostulacion(ctx context.Context, postulanteID int, cedula string, req *GuardarPostulacionRequest) error
	EnviarPostulacion(ctx context.Context, postulanteID int, cedula string, aceptaConsentimiento bool) (string, error)
}

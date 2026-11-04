package domain

import "context"

type PostulanteRepository interface {
	Create(ctx context.Context, postulante *Postulante) error
	FindByCedula(ctx context.Context, cedula string) (*Postulante, error)
	FindByCorreo(ctx context.Context, correo string) (*Postulante, error)
	ExistsByCedula(ctx context.Context, cedula string) (bool, error)
	ExistsByCorreo(ctx context.Context, correo string) (bool, error)
	UpdatePassword(ctx context.Context, id int, passwordHash string) error
}

type RegistroCivilClient interface {
	VerificarIdentidad(ctx context.Context, cedula string, fechaExpedicion string, codigoDactilar string) (*RegistroCivilData, error)
}

type IntentosLimiter interface {
	IncrementarIntento(ctx context.Context, cedula string) (int, error)
	VerificarBloqueo(ctx context.Context, cedula string) (bool, int, error)
	ReiniciarIntentos(ctx context.Context, cedula string) error
	EstablecerBloqueo(ctx context.Context, cedula string, ttl int) error
}

type TokenManager interface {
	GenerarToken(postulanteID int, cedula string) (string, error)
	ValidarToken(token string) (int, string, error)
}

type CatalogoRepository interface {
	ObtenerPaises(ctx context.Context) ([]Pais, error)
	ObtenerProvincias(ctx context.Context) ([]Provincia, error)
	ObtenerCantonesPorProvincia(ctx context.Context, provinciaID int) ([]Canton, error)
	ExisteProvincia(ctx context.Context, provinciaID int) (bool, error)
}

type CatalogoCache interface {
	ObtenerPaises(ctx context.Context) ([]Pais, error)
	GuardarPaises(ctx context.Context, paises []Pais, ttl int) error
	ObtenerProvincias(ctx context.Context) ([]Provincia, error)
	GuardarProvincias(ctx context.Context, provincias []Provincia, ttl int) error
	ObtenerCantonesPorProvincia(ctx context.Context, provinciaID int) ([]Canton, error)
	GuardarCantonesPorProvincia(ctx context.Context, provinciaID int, cantones []Canton, ttl int) error
	InvalidarCatalogos(ctx context.Context) error
}

type HistorialLaboralClient interface {
	VerificarSector(ctx context.Context, cedula string, sector string) (bool, string, error)
}

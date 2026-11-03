package domain

import "context"

type PostulanteRepository interface {
	Create(ctx context.Context, postulante *Postulante) error
	FindByCedula(ctx context.Context, cedula string) (*Postulante, error)
	FindByCorreo(ctx context.Context, correo string) (*Postulante, error)
	ExistsByCedula(ctx context.Context, cedula string) (bool, error)
	ExistsByCorreo(ctx context.Context, correo string) (bool, error)
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

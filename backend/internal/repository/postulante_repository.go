package repository

import (
	"context"
	"errors"
	"time"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/postgres"
	"github.com/jackc/pgx/v5"
)

type PostulanteRepository struct {
	pg *postgres.Postgres
}

func NewPostulanteRepository(pg *postgres.Postgres) *PostulanteRepository {
	return &PostulanteRepository{pg: pg}
}

func (r *PostulanteRepository) Create(ctx context.Context, postulante *domain.Postulante) error {
	query := `
		INSERT INTO postulantes (cedula, correo, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	postulante.CreatedAt = time.Now()

	err := r.pg.Pool().QueryRow(ctx, query,
		postulante.Cedula,
		postulante.Correo,
		postulante.PasswordHash,
		postulante.CreatedAt,
	).Scan(&postulante.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostulanteRepository) FindByCedula(ctx context.Context, cedula string) (*domain.Postulante, error) {
	query := `
		SELECT id, cedula, correo, password_hash, created_at
		FROM postulantes
		WHERE cedula = $1
	`

	var postulante domain.Postulante
	err := r.pg.Pool().QueryRow(ctx, query, cedula).Scan(
		&postulante.ID,
		&postulante.Cedula,
		&postulante.Correo,
		&postulante.PasswordHash,
		&postulante.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &postulante, nil
}

func (r *PostulanteRepository) FindByCorreo(ctx context.Context, correo string) (*domain.Postulante, error) {
	query := `
		SELECT id, cedula, correo, password_hash, created_at
		FROM postulantes
		WHERE correo = $1
	`

	var postulante domain.Postulante
	err := r.pg.Pool().QueryRow(ctx, query, correo).Scan(
		&postulante.ID,
		&postulante.Cedula,
		&postulante.Correo,
		&postulante.PasswordHash,
		&postulante.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &postulante, nil
}

func (r *PostulanteRepository) ExistsByCedula(ctx context.Context, cedula string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM postulantes WHERE cedula = $1)`

	var exists bool
	err := r.pg.Pool().QueryRow(ctx, query, cedula).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PostulanteRepository) ExistsByCorreo(ctx context.Context, correo string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM postulantes WHERE correo = $1)`

	var exists bool
	err := r.pg.Pool().QueryRow(ctx, query, correo).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

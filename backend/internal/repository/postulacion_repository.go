package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostulacionRepository struct {
	pool *pgxpool.Pool
}

func NewPostulacionRepository(pool *pgxpool.Pool) *PostulacionRepository {
	return &PostulacionRepository{pool: pool}
}

func (r *PostulacionRepository) ObtenerPorPostulanteID(ctx context.Context, postulanteID int) (*domain.Postulacion, error) {
	query := `SELECT id, postulante_id, codigo, sector, estado, COALESCE(motivo_cierre, ''), created_at, updated_at 
	          FROM postulaciones WHERE postulante_id = $1`

	var postulacion domain.Postulacion
	err := r.pool.QueryRow(ctx, query, postulanteID).Scan(
		&postulacion.ID,
		&postulacion.PostulanteID,
		&postulacion.Codigo,
		&postulacion.Sector,
		&postulacion.Estado,
		&postulacion.MotivoCierre,
		&postulacion.CreatedAt,
		&postulacion.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error al obtener postulación: %w", err)
	}

	return &postulacion, nil
}

func (r *PostulacionRepository) Crear(ctx context.Context, tx pgx.Tx, postulacion *domain.Postulacion) error {
	query := `INSERT INTO postulaciones (postulante_id, codigo, sector, estado, motivo_cierre)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	err := tx.QueryRow(ctx, query,
		postulacion.PostulanteID,
		postulacion.Codigo,
		postulacion.Sector,
		postulacion.Estado,
		postulacion.MotivoCierre,
	).Scan(&postulacion.ID, &postulacion.CreatedAt, &postulacion.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error al crear postulación: %w", err)
	}

	return nil
}

func (r *PostulacionRepository) CrearSinTx(ctx context.Context, postulacion *domain.Postulacion) error {
	query := `INSERT INTO postulaciones (postulante_id, codigo, sector, estado, motivo_cierre)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		postulacion.PostulanteID,
		postulacion.Codigo,
		postulacion.Sector,
		postulacion.Estado,
		postulacion.MotivoCierre,
	).Scan(&postulacion.ID, &postulacion.CreatedAt, &postulacion.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error al crear postulación: %w", err)
	}

	return nil
}

func (r *PostulacionRepository) Actualizar(ctx context.Context, tx pgx.Tx, postulacion *domain.Postulacion) error {
	query := `UPDATE postulaciones 
	          SET sector = $2, estado = $3, motivo_cierre = $4, updated_at = CURRENT_TIMESTAMP
	          WHERE id = $1`

	_, err := tx.Exec(ctx, query,
		postulacion.ID,
		postulacion.Sector,
		postulacion.Estado,
		postulacion.MotivoCierre,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar postulación: %w", err)
	}

	return nil
}

func (r *PostulacionRepository) GenerarCodigoUnico(ctx context.Context) (string, error) {
	// El UUID evita colisiones entre postulantes concurrentes y cumple el límite VARCHAR(50).
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("error generando identificador de postulación: %w", err)
	}

	return fmt.Sprintf("POST-%s-%s", time.Now().Format("20060102"), id.String()), nil
}

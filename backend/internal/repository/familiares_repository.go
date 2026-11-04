package repository

import (
	"context"
	"fmt"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FamiliaresRepository struct {
	pool *pgxpool.Pool
}

func NewFamiliaresRepository(pool *pgxpool.Pool) *FamiliaresRepository {
	return &FamiliaresRepository{pool: pool}
}

func (r *FamiliaresRepository) EliminarPorPostulacionID(ctx context.Context, tx pgx.Tx, postulacionID int) error {
	query := `DELETE FROM familiares_iess WHERE postulacion_id = $1`
	_, err := tx.Exec(ctx, query, postulacionID)
	if err != nil {
		return fmt.Errorf("error al eliminar familiares: %w", err)
	}
	return nil
}

func (r *FamiliaresRepository) Guardar(ctx context.Context, tx pgx.Tx, familiares []domain.FamiliarIESS) error {
	if len(familiares) == 0 {
		return nil
	}

	query := `INSERT INTO familiares_iess (postulacion_id, nombres, parentesco, institucion, area, cargo)
	          VALUES ($1, $2, $3, $4, $5, $6)`

	for _, familiar := range familiares {
		_, err := tx.Exec(ctx, query,
			familiar.PostulacionID,
			familiar.Nombres,
			familiar.Parentesco,
			familiar.Institucion,
			familiar.Area,
			familiar.Cargo,
		)
		if err != nil {
			return fmt.Errorf("error al guardar familiar: %w", err)
		}
	}

	return nil
}

func (r *FamiliaresRepository) ObtenerPorPostulacionID(ctx context.Context, postulacionID int) ([]domain.FamiliarIESS, error) {
	query := `SELECT id, postulacion_id, nombres, parentesco, institucion, area, cargo, created_at
	          FROM familiares_iess WHERE postulacion_id = $1`

	rows, err := r.pool.Query(ctx, query, postulacionID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener familiares: %w", err)
	}
	defer rows.Close()

	var familiares []domain.FamiliarIESS
	for rows.Next() {
		var familiar domain.FamiliarIESS
		err := rows.Scan(
			&familiar.ID,
			&familiar.PostulacionID,
			&familiar.Nombres,
			&familiar.Parentesco,
			&familiar.Institucion,
			&familiar.Area,
			&familiar.Cargo,
			&familiar.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando familiar: %w", err)
		}
		familiares = append(familiares, familiar)
	}

	return familiares, nil
}

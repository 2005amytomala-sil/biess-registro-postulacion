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

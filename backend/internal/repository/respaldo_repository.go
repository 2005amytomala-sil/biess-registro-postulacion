package repository

import (
	"context"
	"fmt"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RespaldoRepository struct {
	pool *pgxpool.Pool
}

func NewRespaldoRepository(pool *pgxpool.Pool) *RespaldoRepository {
	return &RespaldoRepository{pool: pool}
}

func (r *RespaldoRepository) Guardar(ctx context.Context, tx pgx.Tx, respaldo *domain.RespaldoOrganizacion) error {
	query := `INSERT INTO respaldo_organizacion (postulacion_id, tiene_respaldo, nombre_organizacion)
	          VALUES ($1, $2, $3)
	          ON CONFLICT (postulacion_id) DO UPDATE SET
	          tiene_respaldo = EXCLUDED.tiene_respaldo,
	          nombre_organizacion = EXCLUDED.nombre_organizacion`

	_, err := tx.Exec(ctx, query,
		respaldo.PostulacionID,
		respaldo.TieneRespaldo,
		respaldo.NombreOrganizacion,
	)

	if err != nil {
		return fmt.Errorf("error al guardar respaldo: %w", err)
	}

	return nil
}

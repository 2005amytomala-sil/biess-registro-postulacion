package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RespaldoRepository struct {
	pool *pgxpool.Pool
}

func NewRespaldoRepository(pool *pgxpool.Pool) *RespaldoRepository {
	return &RespaldoRepository{pool: pool}
}

func (r *RespaldoRepository) Guardar(ctx context.Context, respaldo *domain.RespaldoOrganizacion) error {
	query := `INSERT INTO respaldo_organizacion (postulacion_id, tiene_respaldo, nombre_organizacion)
	          VALUES ($1, $2, $3)
	          ON CONFLICT (postulacion_id) DO UPDATE SET
	          tiene_respaldo = EXCLUDED.tiene_respaldo,
	          nombre_organizacion = EXCLUDED.nombre_organizacion`

	_, err := postgres.Conn(ctx, r.pool).Exec(ctx, query,
		respaldo.PostulacionID,
		respaldo.TieneRespaldo,
		respaldo.NombreOrganizacion,
	)

	if err != nil {
		return fmt.Errorf("error al guardar respaldo: %w", err)
	}

	return nil
}

func (r *RespaldoRepository) ObtenerPorPostulacionID(ctx context.Context, postulacionID int) (*domain.RespaldoOrganizacion, error) {
	query := `SELECT postulacion_id, tiene_respaldo, nombre_organizacion, created_at
	          FROM respaldo_organizacion WHERE postulacion_id = $1`

	var respaldo domain.RespaldoOrganizacion
	err := postgres.Conn(ctx, r.pool).QueryRow(ctx, query, postulacionID).Scan(
		&respaldo.PostulacionID,
		&respaldo.TieneRespaldo,
		&respaldo.NombreOrganizacion,
		&respaldo.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener respaldo: %w", err)
	}

	return &respaldo, nil
}

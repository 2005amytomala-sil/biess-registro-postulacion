package repository

import (
	"context"
	"fmt"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FormacionRepository struct {
	pool *pgxpool.Pool
}

func NewFormacionRepository(pool *pgxpool.Pool) *FormacionRepository {
	return &FormacionRepository{pool: pool}
}

func (r *FormacionRepository) EliminarPorPostulacionID(ctx context.Context, tx pgx.Tx, postulacionID int) error {
	query := `DELETE FROM formacion_academica WHERE postulacion_id = $1`
	_, err := tx.Exec(ctx, query, postulacionID)
	if err != nil {
		return fmt.Errorf("error al eliminar formación: %w", err)
	}
	return nil
}

func (r *FormacionRepository) Guardar(ctx context.Context, tx pgx.Tx, formacion []domain.FormacionAcademica) error {
	if len(formacion) == 0 {
		return nil
	}

	query := `INSERT INTO formacion_academica (postulacion_id, nivel, institucion, tipo, titulo, area_estudios, registro_senescyt, fecha_registro_senescyt)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	for _, f := range formacion {
		_, err := tx.Exec(ctx, query,
			f.PostulacionID,
			f.Nivel,
			f.Institucion,
			f.Tipo,
			f.Titulo,
			f.AreaEstudios,
			f.RegistroSenescyt,
			f.FechaRegistroSenescyt,
		)
		if err != nil {
			return fmt.Errorf("error al guardar formación: %w", err)
		}
	}

	return nil
}

func (r *FormacionRepository) ObtenerPorPostulacionID(ctx context.Context, postulacionID int) ([]domain.FormacionAcademica, error) {
	query := `SELECT id, postulacion_id, nivel, institucion, tipo, titulo, area_estudios, registro_senescyt, fecha_registro_senescyt, created_at
	          FROM formacion_academica WHERE postulacion_id = $1`

	rows, err := r.pool.Query(ctx, query, postulacionID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener formación: %w", err)
	}
	defer rows.Close()

	var formaciones []domain.FormacionAcademica
	for rows.Next() {
		var formacion domain.FormacionAcademica
		err := rows.Scan(
			&formacion.ID,
			&formacion.PostulacionID,
			&formacion.Nivel,
			&formacion.Institucion,
			&formacion.Tipo,
			&formacion.Titulo,
			&formacion.AreaEstudios,
			&formacion.RegistroSenescyt,
			&formacion.FechaRegistroSenescyt,
			&formacion.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando formación: %w", err)
		}
		formaciones = append(formaciones, formacion)
	}

	return formaciones, nil
}

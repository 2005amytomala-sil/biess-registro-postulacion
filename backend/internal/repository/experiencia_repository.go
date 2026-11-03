package repository

import (
	"context"
	"fmt"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExperienciaRepository struct {
	pool *pgxpool.Pool
}

func NewExperienciaRepository(pool *pgxpool.Pool) *ExperienciaRepository {
	return &ExperienciaRepository{pool: pool}
}

func (r *ExperienciaRepository) EliminarPorPostulacionID(ctx context.Context, tx pgx.Tx, postulacionID int) error {
	query := `DELETE FROM experiencia_profesional WHERE postulacion_id = $1`
	_, err := tx.Exec(ctx, query, postulacionID)
	if err != nil {
		return fmt.Errorf("error al eliminar experiencia: %w", err)
	}
	return nil
}

func (r *ExperienciaRepository) Guardar(ctx context.Context, tx pgx.Tx, experiencia []domain.ExperienciaProfesional) error {
	if len(experiencia) == 0 {
		return nil
	}

	query := `INSERT INTO experiencia_profesional (postulacion_id, tipo, cargo, institucion, fecha_inicio, fecha_fin, actividades_relacionadas, descripcion)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	for _, exp := range experiencia {
		_, err := tx.Exec(ctx, query,
			exp.PostulacionID,
			exp.Tipo,
			exp.Cargo,
			exp.Institucion,
			exp.FechaInicio,
			exp.FechaFin,
			exp.ActividadesRelacionadas,
			exp.Descripcion,
		)
		if err != nil {
			return fmt.Errorf("error al guardar experiencia: %w", err)
		}
	}

	return nil
}

package repository

import (
	"context"
	"fmt"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExperienciaRepository struct {
	pool *pgxpool.Pool
}

func NewExperienciaRepository(pool *pgxpool.Pool) *ExperienciaRepository {
	return &ExperienciaRepository{pool: pool}
}

func (r *ExperienciaRepository) EliminarPorPostulacionID(ctx context.Context, postulacionID int) error {
	query := `DELETE FROM experiencia_profesional WHERE postulacion_id = $1`
	_, err := postgres.Conn(ctx, r.pool).Exec(ctx, query, postulacionID)
	if err != nil {
		return fmt.Errorf("error al eliminar experiencia: %w", err)
	}
	return nil
}

func (r *ExperienciaRepository) Guardar(ctx context.Context, experiencia []domain.ExperienciaProfesional) error {
	if len(experiencia) == 0 {
		return nil
	}

	query := `INSERT INTO experiencia_profesional (postulacion_id, tipo, cargo, institucion, fecha_inicio, fecha_fin, actividades_relacionadas, descripcion)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	for _, exp := range experiencia {
		_, err := postgres.Conn(ctx, r.pool).Exec(ctx, query,
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

func (r *ExperienciaRepository) ObtenerPorPostulacionID(ctx context.Context, postulacionID int) ([]domain.ExperienciaProfesional, error) {
	query := `SELECT id, postulacion_id, tipo, cargo, institucion, fecha_inicio, fecha_fin, actividades_relacionadas, descripcion, created_at
	          FROM experiencia_profesional WHERE postulacion_id = $1`

	rows, err := postgres.Conn(ctx, r.pool).Query(ctx, query, postulacionID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener experiencia: %w", err)
	}
	defer rows.Close()

	var experiencias []domain.ExperienciaProfesional
	for rows.Next() {
		var experiencia domain.ExperienciaProfesional
		err := rows.Scan(
			&experiencia.ID,
			&experiencia.PostulacionID,
			&experiencia.Tipo,
			&experiencia.Cargo,
			&experiencia.Institucion,
			&experiencia.FechaInicio,
			&experiencia.FechaFin,
			&experiencia.ActividadesRelacionadas,
			&experiencia.Descripcion,
			&experiencia.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando experiencia: %w", err)
		}
		experiencias = append(experiencias, experiencia)
	}

	return experiencias, nil
}

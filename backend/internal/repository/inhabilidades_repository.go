package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InhabilidadesRepository struct {
	pool *pgxpool.Pool
}

func NewInhabilidadesRepository(pool *pgxpool.Pool) *InhabilidadesRepository {
	return &InhabilidadesRepository{pool: pool}
}

func (r *InhabilidadesRepository) Guardar(ctx context.Context, tx pgx.Tx, inhabilidades []domain.Inhabilidad) error {
	if len(inhabilidades) == 0 {
		return nil
	}

	query := `INSERT INTO inhabilidades (postulacion_id, codigo_pregunta, respuesta, detalle)
	          VALUES ($1, $2, $3, $4)
	          ON CONFLICT (postulacion_id, codigo_pregunta) DO UPDATE SET
	          respuesta = EXCLUDED.respuesta,
	          detalle = EXCLUDED.detalle`

	for _, inh := range inhabilidades {
		var detalleJSONB interface{}
		if inh.Detalle != "" {
			// Convertir string a JSONB
			var jsonMap map[string]interface{}
			if err := json.Unmarshal([]byte(inh.Detalle), &jsonMap); err == nil {
				detalleJSONB = jsonMap
			} else {
				detalleJSONB = inh.Detalle
			}
		}

		_, err := tx.Exec(ctx, query,
			inh.PostulacionID,
			inh.CodigoPregunta,
			inh.Respuesta,
			detalleJSONB,
		)
		if err != nil {
			return fmt.Errorf("error al guardar inhabilidad: %w", err)
		}
	}

	return nil
}

func (r *InhabilidadesRepository) ObtenerPorPostulacionID(ctx context.Context, postulacionID int) ([]domain.Inhabilidad, error) {
	query := `SELECT postulacion_id, codigo_pregunta, respuesta, detalle, created_at
	          FROM inhabilidades WHERE postulacion_id = $1`

	rows, err := r.pool.Query(ctx, query, postulacionID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener inhabilidades: %w", err)
	}
	defer rows.Close()

	var inhabilidades []domain.Inhabilidad
	for rows.Next() {
		var inhabilidad domain.Inhabilidad
		var detalleJSONB interface{}
		err := rows.Scan(
			&inhabilidad.PostulacionID,
			&inhabilidad.CodigoPregunta,
			&inhabilidad.Respuesta,
			&detalleJSONB,
			&inhabilidad.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando inhabilidad: %w", err)
		}

		// Convertir JSONB a string si es necesario
		if detalleJSONB != nil {
			if jsonBytes, err := json.Marshal(detalleJSONB); err == nil {
				inhabilidad.Detalle = string(jsonBytes)
			} else {
				inhabilidad.Detalle = ""
			}
		} else {
			inhabilidad.Detalle = ""
		}

		inhabilidades = append(inhabilidades, inhabilidad)
	}

	return inhabilidades, nil
}

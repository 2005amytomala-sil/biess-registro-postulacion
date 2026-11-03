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

package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/postgres"
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

const selectPostulacion = `SELECT id, postulante_id, codigo, sector, estado, COALESCE(motivo_cierre, ''), created_at, updated_at 
	          FROM postulaciones WHERE postulante_id = $1`

func (r *PostulacionRepository) obtener(ctx context.Context, query string, postulanteID int) (*domain.Postulacion, error) {
	var postulacion domain.Postulacion
	err := postgres.Conn(ctx, r.pool).QueryRow(ctx, query, postulanteID).Scan(
		&postulacion.ID,
		&postulacion.PostulanteID,
		&postulacion.Codigo,
		&postulacion.Sector,
		&postulacion.Estado,
		&postulacion.MotivoCierre,
		&postulacion.CreatedAt,
		&postulacion.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("postulación del postulante %d: %w", postulanteID, domain.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener postulación: %w", err)
	}
	return &postulacion, nil
}

func (r *PostulacionRepository) ObtenerPorPostulanteID(ctx context.Context, postulanteID int) (*domain.Postulacion, error) {
	return r.obtener(ctx, selectPostulacion, postulanteID)
}

// ObtenerPorPostulanteIDParaActualizar bloquea la fila (FOR UPDATE) dentro de la transacción del contexto
func (r *PostulacionRepository) ObtenerPorPostulanteIDParaActualizar(ctx context.Context, postulanteID int) (*domain.Postulacion, error) {
	return r.obtener(ctx, selectPostulacion+" FOR UPDATE", postulanteID)
}

// ResumenCompletitud consulta, dentro de la transacción, qué secciones obligatorias están persistidas
func (r *PostulacionRepository) ResumenCompletitud(ctx context.Context, postulacionID int) (*domain.ResumenCompletitud, error) {
	var resumen domain.ResumenCompletitud
	query := `SELECT
	            EXISTS(SELECT 1 FROM datos_personales WHERE postulacion_id = $1),
	            EXISTS(SELECT 1 FROM domicilio_contacto WHERE postulacion_id = $1),
	            EXISTS(SELECT 1 FROM respaldo_organizacion WHERE postulacion_id = $1),
	            COALESCE((SELECT tiene_respaldo FROM respaldo_organizacion WHERE postulacion_id = $1), FALSE),
	            COALESCE((SELECT nombre_organizacion FROM respaldo_organizacion WHERE postulacion_id = $1), ''),
	            COALESCE((SELECT array_agg(codigo_pregunta) FROM inhabilidades WHERE postulacion_id = $1), '{}')`
	err := postgres.Conn(ctx, r.pool).QueryRow(ctx, query, postulacionID).Scan(
		&resumen.DatosPersonales,
		&resumen.Domicilio,
		&resumen.Respaldo,
		&resumen.TieneRespaldo,
		&resumen.NombreOrganizacion,
		&resumen.InhabilidadesRespondidas,
	)
	if err != nil {
		return nil, fmt.Errorf("error al consultar completitud de la postulación: %w", err)
	}
	return &resumen, nil
}

func (r *PostulacionRepository) Crear(ctx context.Context, postulacion *domain.Postulacion) error {
	query := `INSERT INTO postulaciones (postulante_id, codigo, sector, estado, motivo_cierre)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	err := postgres.Conn(ctx, r.pool).QueryRow(ctx, query,
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

func (r *PostulacionRepository) Actualizar(ctx context.Context, postulacion *domain.Postulacion) error {
	query := `UPDATE postulaciones 
	          SET sector = $2, estado = $3, motivo_cierre = NULLIF($4, ''), codigo = $5, updated_at = CURRENT_TIMESTAMP
	          WHERE id = $1`

	_, err := postgres.Conn(ctx, r.pool).Exec(ctx, query,
		postulacion.ID,
		postulacion.Sector,
		postulacion.Estado,
		postulacion.MotivoCierre,
		postulacion.Codigo,
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

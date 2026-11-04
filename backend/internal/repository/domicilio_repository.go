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

type DomicilioRepository struct {
	pool *pgxpool.Pool
}

func NewDomicilioRepository(pool *pgxpool.Pool) *DomicilioRepository {
	return &DomicilioRepository{pool: pool}
}

func (r *DomicilioRepository) Guardar(ctx context.Context, domicilio *domain.DomicilioContacto) error {
	query := `INSERT INTO domicilio_contacto 
	          (postulacion_id, pais_id, provincia_id, canton_id, provincia_texto, ciudad_texto, celular, correo_1, correo_2, 
	           calle_principal, numero, calle_secundaria, sector, referencia, telefono_domicilio)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	          ON CONFLICT (postulacion_id) DO UPDATE SET
	          pais_id = EXCLUDED.pais_id,
	          provincia_id = EXCLUDED.provincia_id,
	          canton_id = EXCLUDED.canton_id,
	          provincia_texto = EXCLUDED.provincia_texto,
	          ciudad_texto = EXCLUDED.ciudad_texto,
	          celular = EXCLUDED.celular,
	          correo_1 = EXCLUDED.correo_1,
	          correo_2 = EXCLUDED.correo_2,
	          calle_principal = EXCLUDED.calle_principal,
	          numero = EXCLUDED.numero,
	          calle_secundaria = EXCLUDED.calle_secundaria,
	          sector = EXCLUDED.sector,
	          referencia = EXCLUDED.referencia,
	          telefono_domicilio = EXCLUDED.telefono_domicilio`

	_, err := postgres.Conn(ctx, r.pool).Exec(ctx, query,
		domicilio.PostulacionID,
		domicilio.PaisID,
		domicilio.ProvinciaID,
		domicilio.CantonID,
		domicilio.ProvinciaTexto,
		domicilio.CiudadTexto,
		domicilio.Celular,
		domicilio.Correo1,
		domicilio.Correo2,
		domicilio.CallePrincipal,
		domicilio.Numero,
		domicilio.CalleSecundaria,
		domicilio.Sector,
		domicilio.Referencia,
		domicilio.TelefonoDomicilio,
	)

	if err != nil {
		return fmt.Errorf("error al guardar domicilio: %w", err)
	}

	return nil
}

func (r *DomicilioRepository) ObtenerPorPostulacionID(ctx context.Context, postulacionID int) (*domain.DomicilioContacto, error) {
	query := `SELECT postulacion_id, pais_id, provincia_id, canton_id, provincia_texto, ciudad_texto, celular, correo_1, correo_2, 
	          calle_principal, numero, calle_secundaria, sector, referencia, telefono_domicilio, created_at
	          FROM domicilio_contacto WHERE postulacion_id = $1`

	var domicilio domain.DomicilioContacto
	err := postgres.Conn(ctx, r.pool).QueryRow(ctx, query, postulacionID).Scan(
		&domicilio.PostulacionID,
		&domicilio.PaisID,
		&domicilio.ProvinciaID,
		&domicilio.CantonID,
		&domicilio.ProvinciaTexto,
		&domicilio.CiudadTexto,
		&domicilio.Celular,
		&domicilio.Correo1,
		&domicilio.Correo2,
		&domicilio.CallePrincipal,
		&domicilio.Numero,
		&domicilio.CalleSecundaria,
		&domicilio.Sector,
		&domicilio.Referencia,
		&domicilio.TelefonoDomicilio,
		&domicilio.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener domicilio: %w", err)
	}

	return &domicilio, nil
}

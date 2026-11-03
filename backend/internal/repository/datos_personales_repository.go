package repository

import (
	"context"
	"fmt"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatosPersonalesRepository struct {
	pool *pgxpool.Pool
}

func NewDatosPersonalesRepository(pool *pgxpool.Pool) *DatosPersonalesRepository {
	return &DatosPersonalesRepository{pool: pool}
}

func (r *DatosPersonalesRepository) Guardar(ctx context.Context, tx pgx.Tx, datos *domain.DatosPersonales) error {
	query := `INSERT INTO datos_personales 
	          (postulacion_id, nombres, apellidos, cedula, genero, fecha_nacimiento, nacionalidad, estado_civil, pareja_nombres, pareja_cedula)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	          ON CONFLICT (postulacion_id) DO UPDATE SET
	          nombres = EXCLUDED.nombres,
	          apellidos = EXCLUDED.apellidos,
	          cedula = EXCLUDED.cedula,
	          genero = EXCLUDED.genero,
	          fecha_nacimiento = EXCLUDED.fecha_nacimiento,
	          nacionalidad = EXCLUDED.nacionalidad,
	          estado_civil = EXCLUDED.estado_civil,
	          pareja_nombres = EXCLUDED.pareja_nombres,
	          pareja_cedula = EXCLUDED.pareja_cedula`

	_, err := tx.Exec(ctx, query,
		datos.PostulacionID,
		datos.Nombres,
		datos.Apellidos,
		datos.Cedula,
		datos.Genero,
		datos.FechaNacimiento,
		datos.Nacionalidad,
		datos.EstadoCivil,
		datos.ParejaNombres,
		datos.ParejaCedula,
	)

	if err != nil {
		return fmt.Errorf("error al guardar datos personales: %w", err)
	}

	return nil
}

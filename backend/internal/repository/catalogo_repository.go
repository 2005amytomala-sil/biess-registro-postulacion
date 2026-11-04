package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogoRepository struct {
	pool *pgxpool.Pool
}

func NewCatalogoRepository(pool *pgxpool.Pool) *CatalogoRepository {
	return &CatalogoRepository{pool: pool}
}

func (r *CatalogoRepository) ObtenerPaises(ctx context.Context) ([]domain.Pais, error) {
	query := `SELECT id, nombre, created_at FROM paises ORDER BY nombre`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar paises: %w", err)
	}
	defer rows.Close()

	var paises []domain.Pais
	for rows.Next() {
		var pais domain.Pais
		if err := rows.Scan(&pais.ID, &pais.Nombre, &pais.CreatedAt); err != nil {
			return nil, fmt.Errorf("error al escanear pais: %w", err)
		}
		paises = append(paises, pais)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar paises: %w", err)
	}

	return paises, nil
}

func (r *CatalogoRepository) ObtenerProvincias(ctx context.Context) ([]domain.Provincia, error) {
	query := `SELECT id, nombre, created_at FROM provincias ORDER BY id`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar provincias: %w", err)
	}
	defer rows.Close()

	var provincias []domain.Provincia
	for rows.Next() {
		var provincia domain.Provincia
		if err := rows.Scan(&provincia.ID, &provincia.Nombre, &provincia.CreatedAt); err != nil {
			return nil, fmt.Errorf("error al escanear provincia: %w", err)
		}
		provincias = append(provincias, provincia)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar provincias: %w", err)
	}

	return provincias, nil
}

func (r *CatalogoRepository) ObtenerCantonesPorProvincia(ctx context.Context, provinciaID int) ([]domain.Canton, error) {
	query := `SELECT id, provincia_id, nombre, created_at FROM cantones WHERE provincia_id = $1 ORDER BY nombre`

	rows, err := r.pool.Query(ctx, query, provinciaID)
	if err != nil {
		return nil, fmt.Errorf("error al consultar cantones: %w", err)
	}
	defer rows.Close()

	var cantones []domain.Canton
	for rows.Next() {
		var canton domain.Canton
		if err := rows.Scan(&canton.ID, &canton.ProvinciaID, &canton.Nombre, &canton.CreatedAt); err != nil {
			return nil, fmt.Errorf("error al escanear canton: %w", err)
		}
		cantones = append(cantones, canton)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar cantones: %w", err)
	}

	return cantones, nil
}

func (r *CatalogoRepository) ObtenerPaisPorID(ctx context.Context, paisID int) (*domain.Pais, error) {
	query := `SELECT id, nombre, created_at FROM paises WHERE id = $1`

	var pais domain.Pais
	err := r.pool.QueryRow(ctx, query, paisID).Scan(&pais.ID, &pais.Nombre, &pais.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener país: %w", err)
	}
	return &pais, nil
}

func (r *CatalogoRepository) ExisteProvincia(ctx context.Context, provinciaID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM provincias WHERE id = $1)`

	var existe bool
	err := r.pool.QueryRow(ctx, query, provinciaID).Scan(&existe)
	if err != nil {
		return false, fmt.Errorf("error al verificar existencia de provincia: %w", err)
	}

	return existe, nil
}

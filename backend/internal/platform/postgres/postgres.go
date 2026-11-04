package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.PostgresConfig) (*Postgres, error) {
	dsn := cfg.DSN()
	if cfg.Host == "localhost" {
		dsn = fmt.Sprintf("host=127.0.0.1 port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	}
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("error parsing postgres config: %w", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging postgres: %w", err)
	}

	return &Postgres{pool: pool}, nil
}

func (p *Postgres) Pool() *pgxpool.Pool {
	return p.pool
}

func (p *Postgres) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

type txKey struct{}

// DBTX es el subconjunto común de pgxpool.Pool y pgx.Tx que usan los repositorios
type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Conn devuelve la transacción activa en el contexto o, si no hay, el pool
func Conn(ctx context.Context, pool *pgxpool.Pool) DBTX {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return pool
}

// WithTx implementa domain.TxManager: BEGIN al inicio, COMMIT si fn no falla, ROLLBACK en caso contrario.
// La transacción viaja en el contexto; si ya hay una activa, fn se ejecuta dentro de ella.
func (p *Postgres) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return fn(ctx)
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p) // re-throw el panic después del rollback
		}
	}()

	if err := fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("error en rollback: %v (error original: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error al hacer commit: %w", err)
	}

	return nil
}

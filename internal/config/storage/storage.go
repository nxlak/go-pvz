package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Pool *pgxpool.Pool
}

func Connect(connString string, maxConns int) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = int32(maxConns)

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &Storage{Pool: pool}, nil
}

func (s *Storage) Close() {
	if s != nil && s.Pool != nil {
		s.Pool.Close()
	}
}

// USAGE
//db, err := db.Connect(os.Getenv("DATABASE_URL"), MAX_CONN)
//defer db.Close()

//postgres://<USERNAME>:<PASSWORD>@<HOST>:<PORT>/<DBNAME>?sslmode=disable = DATABASE_URL

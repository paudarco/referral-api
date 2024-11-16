package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paudarco/referral-api/internal/config"
)

func NewPostresPool(cfg config.DB) (*pgxpool.Pool, error) {
	configStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&client_encoding=UTF8",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(configStr)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}

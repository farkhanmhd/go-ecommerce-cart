package app

import (
	"context"
	"time"

	"farkhanmhd/go-ecommerce-cart/common/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPGPool(dbUrl string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dbUrl)

	utils.PanicIfError(err)

	config.MaxConns = 20
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	return pool
}

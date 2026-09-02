package db

import (
	"context"
	"embed"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var EmbededMigration embed.FS

func NewPool(connectionUrl string, ctx context.Context) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(ctx, connectionUrl)

	if err != nil {
		log.Println("Error occured while creating a pool.")
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		log.Println("Error: Could not reach the database:", err)
		return nil, err
	}

	log.Println("Connected to Supabase.")

	return pool, nil

}

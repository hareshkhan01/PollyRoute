package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/hareshkhan01/PollyRoute/internals/config"
	"github.com/hareshkhan01/PollyRoute/internals/db"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Failed to Load the config: ", err)
		return
	}

	databaseUrl := cfg.DATABASE_URL
	ctx := context.Background()
	migrateDb, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		fmt.Println("Failed to open a db connection for migration: ", err)
		return
	}

	goose.SetBaseFS(db.EmbededMigration)

	if err := goose.SetDialect("postgres"); err != nil {
		fmt.Println("Failed to SetDialect: ", err)
		return
	}

	command := "up"
	var args []string
	if len(os.Args) > 1 {
		command = os.Args[1]
		args = os.Args[2:]
	}

	if err := goose.RunContext(ctx, command, migrateDb, "migrations", args...); err != nil {
		fmt.Println("Failed Run Goose: ", err)
		return
	}
	fmt.Println("Migration Successfull")
}

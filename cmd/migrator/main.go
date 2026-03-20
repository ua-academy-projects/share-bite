package main

import (
	"context"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

func main() {
	ctx := context.Background()

	if err := config.Load(".env"); err != nil {
		logger.Fatal(ctx, "config load: ", err)
	}

	client, err := pg.NewClient(ctx, config.Config().Postgres.Dsn())
	if err != nil {
		logger.Fatal(ctx, "new database client: ", err)
	}
	if err := client.DB().Ping(ctx); err != nil {
		logger.Fatal(ctx, "ping database: ", err)
	}

	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		logger.Fatal(ctx, "set migrations dialect: ", err)
	}

	db := stdlib.OpenDBFromPool(client.DB().Pool())
	if err := goose.UpContext(ctx, db, config.Config().Postgres.MigrationsDir()); err != nil {
		logger.Fatal(ctx, "migrate up: ", err)
	}

	if err := db.Close(); err != nil {
		logger.Fatal(ctx, "clode db after migrate up: ", err)
	}

	logger.Info(ctx, "migrations applied")
}

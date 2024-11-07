package db

import (
	"database/sql"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/zap"
)

func InitDB(cfg config.Config) *bun.DB {
	connStr := cfg.DBConnString()

	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr)))

	pgdb := bun.NewDB(db, pgdialect.New())

	if cfg.AppEnv() == "DEV" {
		pgdb.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}

	err := pgdb.Ping()

	if err != nil {
		cfg.Logger().Error("Failed to connect to db", zap.Error(err))
		panic("Failed to connect to db")
	}

	cfg.Logger().Info("Connected to database")

	return pgdb
}

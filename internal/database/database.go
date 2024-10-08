package database

import (
	"database/sql"
	"time"

	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func InitDB() *bun.DB {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")), // Load from env
		pgdriver.WithTLSConfig(nil),
		pgdriver.WithUser(os.Getenv("DB_USER")),         // Load from env
		pgdriver.WithPassword(os.Getenv("DB_PASSWORD")), // Load from env
		pgdriver.WithDatabase(os.Getenv("DB_NAME")),     // Load from env
		pgdriver.WithApplicationName("todolist-api"),
		pgdriver.WithTimeout(5*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(5*time.Second),
		pgdriver.WithWriteTimeout(5*time.Second),
	)

	// Create and return a Bun database instance
	return bun.NewDB(sql.OpenDB(pgconn), pgdialect.New())

}

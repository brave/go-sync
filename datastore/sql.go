package datastore

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

const sqlURLEnvKey = "SQL_DATABASE_URL"

// SQLDB is a Datastore wrapper around a SQL-based database.
type SQLDB struct {
	*sqlx.DB
	insertQuery string
}

// NewSQLDB returns a SQLDB client to be used.
func NewSQLDB() (*SQLDB, error) {
	sqlURL := os.Getenv(sqlURLEnvKey)
	if len(sqlURL) == 0 {
		return nil, fmt.Errorf("%s must be defined", sqlURLEnvKey)
	}
	migration, err := migrate.New(
		"file://./migrations",
		sqlURL,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to init migrations: %v", err)
	}
	if err = migration.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return nil, fmt.Errorf("Failed to run migrations: %v", err)
		}
	}

	db, err := sqlx.Connect("postgres", sqlURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to SQL DB: %v", err)
	}

	wrappedDB := SQLDB{db, buildInsertQuery()}
	return &wrappedDB, nil
}

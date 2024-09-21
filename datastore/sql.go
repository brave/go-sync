package datastore

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	// import postgres package for migrations
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	// import pgx so it can be used with sqlx
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	sqlURLEnvKey     = "SQL_DATABASE_URL"
	sqlTestURLEnvKey = "SQL_TEST_DATABASE_URL"
	// Default value is defined here, since the .env file will not be loaded
	// because tests are run in the subdirectories where the tests live
	defaultSQLTestURL = "postgres://sync:password@localhost:5434/testing?sslmode=disable"
	// SQLMigrateUpdateIntervalEnvKey is the env var name used to define the frequency
	// of chunked migration within "get update" requests
	SQLMigrateUpdateIntervalEnvKey = "SQL_MIGRATE_UPDATE_INTERVAL"
	// SQLMigrateChunkSizeEnvKey is the env var name used to define the max migration
	// chunk size
	SQLMigrateChunkSizeEnvKey    = "SQL_MIGRATE_CHUNK_SIZE"
	defaultMigrateUpdateInterval = 4
	defaultMigrateChunkSize      = 100
)

//go:embed migrations/*
var migrationFiles embed.FS

// SQLDB is a Datastore wrapper around a SQL-based database.
type SQLDB struct {
	*sqlx.DB
	insertQuery            string
	variations             *SQLVariations
	migrateIntervalPercent float32
	migrateChunkSize       int
}

// NewSQLDB returns a SQLDB client to be used.
func NewSQLDB(isTesting bool) (*SQLDB, error) {
	variations, err := LoadSQLVariations()
	if err != nil {
		return nil, err
	}

	var envKey string
	if isTesting {
		envKey = sqlTestURLEnvKey
	} else {
		envKey = sqlURLEnvKey
	}

	sqlURL := os.Getenv(envKey)
	if sqlURL == "" {
		if isTesting {
			sqlURL = defaultSQLTestURL
		} else {
			return nil, fmt.Errorf("%s must be defined", envKey)
		}
	}
	iofsDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to load iofs driver for migrations: %w", err)
	}
	migration, err := migrate.NewWithSourceInstance(
		"iofs",
		iofsDriver,
		sqlURL,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to init migrations: %w", err)
	}
	if err = migration.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return nil, fmt.Errorf("Failed to run migrations: %w", err)
		}
	}

	db, err := sqlx.Connect("pgx", sqlURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to SQL DB: %w", err)
	}

	if isTesting {
		variations.Ready = true
	}

	migrateInterval, _ := strconv.Atoi(os.Getenv(SQLMigrateUpdateIntervalEnvKey))
	migrateChunkSize, _ := strconv.Atoi(os.Getenv(SQLMigrateChunkSizeEnvKey))

	if migrateInterval <= 0 {
		migrateInterval = defaultMigrateUpdateInterval
	}
	migrateIntervalPercent := 1 / float32(migrateInterval)
	if migrateChunkSize <= 0 {
		migrateChunkSize = defaultMigrateChunkSize
	}

	wrappedDB := SQLDB{db, buildInsertQuery(), variations, migrateIntervalPercent, migrateChunkSize}
	return &wrappedDB, nil
}

// MigrateIntervalPercent returns the percentage of update requests that will perform
// a chunked migration
func (db *SQLDB) MigrateIntervalPercent() float32 {
	return db.migrateIntervalPercent
}

// MigrateChunkSize returns the max chunk size of migration attempts
func (db *SQLDB) MigrateChunkSize() int {
	return db.migrateChunkSize
}

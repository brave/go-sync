package datastore

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

const sqlURLEnvKey = "SQL_DATABASE_URL"
const sqlMigrateUpdateIntervalEnvKey = "SQL_MIGRATE_UPDATE_INTERVAL"
const sqlMigrateChunkSizeEnvKey = "SQL_MIGRATE_CHUNK_SIZE"

const defaultMigrateUpdateInterval = 4
const defaultMigrateChunkSize = 100

// SQLDB is a Datastore wrapper around a SQL-based database.
type SQLDB struct {
	*sqlx.DB
	insertQuery            string
	variations             *SQLVariations
	migrateIntervalPercent float32
	migrateChunkSize       int
}

// NewSQLDB returns a SQLDB client to be used.
func NewSQLDB() (*SQLDB, error) {
	variations, err := LoadSQLVariations()
	if err != nil {
		return nil, err
	}

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

	db, err := sqlx.Connect("pgx", sqlURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to SQL DB: %v", err)
	}

	migrateInterval, _ := strconv.Atoi(os.Getenv(sqlMigrateUpdateIntervalEnvKey))
	migrateChunkSize, _ := strconv.Atoi(os.Getenv(sqlMigrateChunkSizeEnvKey))

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

func (db *SQLDB) MigrateIntervalPercent() float32 {
	return db.migrateIntervalPercent
}

func (db *SQLDB) MigrateChunkSize() int {
	return db.migrateChunkSize
}

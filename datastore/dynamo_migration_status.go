package datastore

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type MigrationStatus struct {
	ChainID       int64  `db:"chain_id"`
	DataType      int    `db:"data_type"`
	EarliestMtime *int64 `db:"earliest_mtime"`
}

func (sqlDB *SQLDB) GetDynamoMigrationStatuses(tx *sqlx.Tx, chainID int64, dataTypes []int) (dataTypeToStatusMap map[int]*MigrationStatus, err error) {
	dataTypeToStatusMap = make(map[int]*MigrationStatus)

	var statuses []MigrationStatus
	err = tx.Select(&statuses, `
		SELECT chain_id, data_type, earliest_mtime
		FROM dynamo_migration_statuses
		WHERE chain_id = $1 AND data_type IN $2
	`, chainID, dataTypes)

	if err != nil {
		return nil, fmt.Errorf("failed to get dynamo migration status: %w", err)
	}

	for _, status := range statuses {
		dataTypeToStatusMap[status.DataType] = &status
	}

	return dataTypeToStatusMap, nil
}

func (sqlDB *SQLDB) UpdateDynamoMigrationStatuses(tx *sqlx.Tx, statuses []*MigrationStatus) error {
	_, err := tx.NamedExec(`
		INSERT INTO dynamo_migration_statuses (chain_id, data_type, earliest_mtime)
		VALUES (:chain_id, :data_type, :earliest_mtime)
			ON CONFLICT DO UPDATE
			SET earliest_mtime = $3
			WHERE earliest_mtime IS NOT NULL AND earliest_mtime > :earliest_mtime
	`, statuses)
	if err != nil {
		return fmt.Errorf("failed to update dynamo migration statuses: %w", err)
	}

	return nil
}

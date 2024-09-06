package datastore

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const chainIDSelectQuery = "SELECT id FROM chains WHERE client_id = $1 FOR UPDATE"

type ChainRow struct {
	ID *int64
}

type MigrationStatus struct {
	ChainID       int64 `db:"chain_id"`
	DataType      int   `db:"data_type"`
	EarliestMtime int64 `db:"earliest_mtime"`
}

func (sqlDB *SQLDB) InsertSyncEntity(tx *sqlx.Tx, entity *SyncEntity) (bool, error) {
	res, err := tx.NamedExec(`
		INSERT INTO entities (
			id, chain_id, data_type, ctime, mtime, specifics, client_defined_unique_tag,
			server_defined_unique_tag, deleted, folder, version, name, originator_cache_guid,
			originator_client_item_id, parent_id, non_unique_name, unique_position
		) VALUES (
			:id, :chain_id, :data_type, :ctime, :mtime, :specifics, :client_defined_unique_tag,
			:server_defined_unique_tag, :deleted, :folder, :version, :name, :originator_cache_guid,
			:originator_client_item_id, :parent_id, :non_unique_name, :unique_position
		) ON CONFLICT DO NOTHING
	`, entity)
	if err != nil {
		return false, fmt.Errorf("failed to insert entity: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected after insert: %w", err)
	}

	// if rows affected is 0, then there must be a conflict. return true to indicate this condition.
	return rowsAffected == 0, nil
}

func (sqlDB *SQLDB) HasItem(tx *sqlx.Tx, chainId int64, idBytes []byte) (bool, error) {
	var exists bool
	err := tx.QueryRowx("SELECT EXISTS(SELECT 1 FROM entities WHERE chain_id = $1 AND id = $2)", chainId, idBytes).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check existence of item: %w", err)
	}
	return exists, nil
}

func (sqlDB *SQLDB) UpdateDynamoMigrationStatuses(tx *sqlx.Tx, chainID int64, data_type_earliest_mtime_map map[int]int64) error {
	var statuses []MigrationStatus
	for dataType, earliestMtime := range data_type_earliest_mtime_map {
		statuses = append(statuses, MigrationStatus{
			ChainID:       chainID,
			DataType:      dataType,
			EarliestMtime: earliestMtime,
		})
	}

	_, err := tx.NamedExec(`
		INSERT INTO dynamo_migration_statuses (chain_id, data_type, earliest_mtime)
		VALUES (:chain_id, :data_type, :earliest_mtime)
			ON CONFLICT (chain_id, data_type) DO UPDATE
			SET earliest_mtime = $3
			WHERE earliest_mtime IS NOT NULL AND earliest_mtime > :earliest_mtime
	`, statuses)
	if err != nil {
		return fmt.Errorf("failed to update dynamo migration statuses: %w", err)
	}

	return nil
}

func (sqlDB *SQLDB) UpdateSyncEntity(tx *sqlx.Tx, entity *SyncEntity, oldVersion int64) (bool, bool, error) {
	whereClause := "WHERE id = :id AND chain_id = :chain_id"
	if *entity.DataType != HistoryTypeID {
		entity.OldVersion = &oldVersion
		whereClause += " AND version = :old_version"
	}

	var query string

	deleted := entity.Deleted != nil && *entity.Deleted
	if deleted {
		query = `DELETE FROM entities ` + whereClause
	} else {
		var updateFields []string
		if entity.UniquePosition != nil {
			updateFields = append(updateFields, "unique_position = :unique_position")
		}
		if entity.ParentID != nil {
			updateFields = append(updateFields, "parent_id = :parent_id")
		}
		if entity.Name != nil {
			updateFields = append(updateFields, "name = :name")
		}
		if entity.NonUniqueName != nil {
			updateFields = append(updateFields, "non_unique_name = :non_unique_name")
		}
		if entity.Folder != nil {
			updateFields = append(updateFields, "folder = :folder")
		}

		var joinedUpdateFields string
		if len(updateFields) > 0 {
			joinedUpdateFields = ", " + strings.Join(updateFields, ", ")
		}
		query = `
			UPDATE entities
			SET version = :version,
				mtime = :mtime,
				specifics = :specifics
		` + joinedUpdateFields + whereClause
	}

	result, err := tx.NamedExec(query, entity)
	if err != nil {
		return false, false, fmt.Errorf("error updating entity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, false, fmt.Errorf("error getting rows affected after update: %w", err)
	}

	return rowsAffected == 0, deleted, nil
}

func (sqlDB *SQLDB) GetAndLockChainID(tx *sqlx.Tx, clientID *string) (*int64, error) {
	// Get chain ID and lock for updates
	clientIDBytes, err := hex.DecodeString(*clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to decode clientID: %w", err)
	}
	row := ChainRow{}
	if err := tx.Get(&row, chainIDSelectQuery, clientIDBytes); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to get chain id: %w", err)
		}
		_, err := tx.Exec("INSERT INTO chains (client_id) VALUES ($1)", clientIDBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to insert chain: %w", err)
		}

		if err = tx.Get(&row, chainIDSelectQuery, clientIDBytes); err != nil {
			return nil, fmt.Errorf("failed to get chain id: %w", err)
		}
	}

	return row.ID, nil
}

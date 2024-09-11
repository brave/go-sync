package datastore

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

var fieldsToInsert = []string{
	"id", "chain_id", "data_type", "ctime", "mtime", "specifics",
	"deleted", "client_defined_unique_tag", "server_defined_unique_tag", "folder", "version",
	"name", "originator_cache_guid", "originator_client_item_id", "parent_id", "non_unique_name",
	"unique_position",
}

type MigrationStatus struct {
	ChainID       int64 `db:"chain_id"`
	DataType      int   `db:"data_type"`
	EarliestMtime int64 `db:"earliest_mtime"`
}

func buildInsertQuery() string {
	var insertValues []string
	var setValues []string
	for _, field := range fieldsToInsert {
		insertValues = append(insertValues, ":"+field)
		setValues = append(setValues, field+" = EXCLUDED."+field)
	}
	joinedFields := strings.Join(fieldsToInsert, ", ")
	joinedInsertValues := strings.Join(insertValues, ", ")
	joinedSetValues := strings.Join(setValues, ", ")
	// We only want to update an existing row if it was previously deleted.
	// If it was not deleted, then it should be considered a conflict
	return `INSERT INTO entities (` + joinedFields + `) VALUES (` + joinedInsertValues +
		`) ON CONFLICT (chain_id, client_defined_unique_tag) DO UPDATE SET ` +
		joinedSetValues + ` WHERE entities.deleted = true`
}

func (sqlDB *SQLDB) InsertSyncEntity(tx *sqlx.Tx, entity *SyncEntity) (conflict bool, err error) {
	res, err := tx.NamedExec(sqlDB.insertQuery, entity)
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

func (sqlDB *SQLDB) HasItem(tx *sqlx.Tx, chainId int64, clientTag string) (exists bool, err error) {
	err = tx.QueryRowx("SELECT EXISTS(SELECT 1 FROM entities WHERE chain_id = $1 AND client_defined_unique_tag = $2)", chainId, clientTag).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check existence of item: %w", err)
	}
	return exists, nil
}

func (sqlDB *SQLDB) GetDynamoMigrationStatus(chainID int64, dataType int) (*MigrationStatus, error) {
	var status MigrationStatus
	err := sqlDB.Get(&status, `
		SELECT chain_id, data_type, earliest_mtime
		FROM dynamo_migration_statuses
		WHERE chain_id = $1 AND data_type = $2
	`, chainID, dataType)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get dynamo migration status: %w", err)
	}

	return &status, nil
}

func (sqlDB *SQLDB) UpdateDynamoMigrationStatuses(tx *sqlx.Tx, statuses []MigrationStatus) error {
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

func (sqlDB *SQLDB) UpdateSyncEntity(tx *sqlx.Tx, entity *SyncEntity, oldVersion int64) (conflict bool, err error) {
	var idColumn string
	if *entity.DataType == HistoryTypeID {
		idColumn = "client_defined_unique_tag"
	} else {
		idColumn = "id"
	}
	whereClause := " WHERE " + idColumn + " = :id AND chain_id = :chain_id AND deleted = false"
	if *entity.DataType != HistoryTypeID {
		entity.OldVersion = &oldVersion
		whereClause += " AND version = :old_version"
	}

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
	if entity.Deleted != nil {
		updateFields = append(updateFields, "deleted = :deleted")
	}

	var joinedUpdateFields string
	if len(updateFields) > 0 {
		joinedUpdateFields = ", " + strings.Join(updateFields, ", ")
	}
	query := `
		UPDATE entities
		SET version = :version,
			mtime = :mtime,
			specifics = :specifics
	` + joinedUpdateFields + whereClause

	result, err := tx.NamedExec(query, entity)
	if err != nil {
		return false, fmt.Errorf("error updating entity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error getting rows affected after update: %w", err)
	}

	return rowsAffected == 0, nil
}

func (sqlDB *SQLDB) GetAndLockChainID(tx *sqlx.Tx, clientID string) (chainID *int64, err error) {
	// Get chain ID and lock for updates
	clientIDBytes, err := hex.DecodeString(clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to decode clientID: %w", err)
	}

	var id int64

	err = tx.QueryRowx(`
		INSERT INTO chains (client_id, last_usage_time) VALUES ($1, $2)
		ON CONFLICT (client_id)
		DO UPDATE SET last_usage_time = EXCLUDED.last_usage_time
		RETURNING id
	`, clientIDBytes, time.Now()).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert chain: %w", err)
	}

	// Once we have completely migrated over to SQL, we can change this to
	// `FOR UPDATE`, and only lock upon commits. We need to lock for updates
	// as we will be deleting older Dynamo items during update requests, and migrating
	// them over to SQL. If another client in the chain updates during this process,
	// the client may not receive some older items.
	_, err = tx.Exec(`SELECT id FROM chains WHERE id = $1 FOR SHARE`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock on chain: %w", err)
	}

	return &id, nil
}

func (sqlDB *SQLDB) GetUpdatesForType(dataType int, clientToken int64, fetchFolders bool, chainID int64, maxSize int) (hasChangesRemaining bool, entities []SyncEntity, err error) {
	var additionalCondition string
	if !fetchFolders {
		additionalCondition = "AND folder = false "
	}
	query := `SELECT * FROM entities
		WHERE chain_id = $1 AND data_type = $2 AND mtime > $3 ` + additionalCondition + `ORDER BY mtime LIMIT $4`

	if err := sqlDB.Select(&entities, query, chainID, dataType, clientToken, maxSize); err != nil {
		return false, nil, fmt.Errorf("failed to get entity updates: %w", err)
	}
	return len(entities) == maxSize, entities, nil
}

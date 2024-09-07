package datastore

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const chainIDSelectQuery = "SELECT id FROM chains WHERE client_id = $1"

var fieldsToInsert = []string{
	"id", "chain_id", "data_type", "ctime", "mtime", "id_is_uuid", "specifics",
	"deleted", "client_defined_unique_tag", "server_defined_unique_tag", "folder", "version",
	"name", "originator_cache_guid", "originator_client_item_id", "parent_id", "non_unique_name",
	"unique_position",
}

type ChainRow struct {
	ID *int64
}

type MigrationStatus struct {
	ChainID       int64 `db:"chain_id"`
	DataType      int   `db:"data_type"`
	EarliestMtime int64 `db:"earliest_mtime"`
}

type CommonSQLX interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...any) (sql.Result, error)
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

func (sqlDB *SQLDB) InsertSyncEntity(tx *sqlx.Tx, entity *SyncEntity) (bool, error) {
	if entity.Deleted != nil && *entity.Deleted {
		return true, nil
	}
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
			ON CONFLICT DO UPDATE
			SET earliest_mtime = $3
			WHERE earliest_mtime IS NOT NULL AND earliest_mtime > :earliest_mtime
	`, statuses)
	if err != nil {
		return fmt.Errorf("failed to update dynamo migration statuses: %w", err)
	}

	return nil
}

func (sqlDB *SQLDB) UpdateSyncEntity(tx *sqlx.Tx, entity *SyncEntity, oldVersion int64) (bool, bool, error) {
	whereClause := " WHERE id = :id AND chain_id = :chain_id AND deleted = false"
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
		return false, false, fmt.Errorf("error updating entity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, false, fmt.Errorf("error getting rows affected after update: %w", err)
	}

	return rowsAffected == 0, entity.Deleted != nil && *entity.Deleted, nil
}

func (sqlDB *SQLDB) GetChainID(tx *sqlx.Tx, clientID string, acquireUpdateLock bool) (*int64, error) {
	// Get chain ID and lock for updates
	clientIDBytes, err := hex.DecodeString(clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to decode clientID: %w", err)
	}
	row := ChainRow{}

	var lockClause string
	if acquireUpdateLock {
		lockClause = " FOR UPDATE"
	}

	var commonSQLX CommonSQLX
	if tx != nil {
		commonSQLX = tx
	} else {
		commonSQLX = sqlDB
	}

	if err := commonSQLX.Get(&row, chainIDSelectQuery+lockClause, clientIDBytes); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to get chain id: %w", err)
		}
		_, err := commonSQLX.Exec("INSERT INTO chains (client_id) VALUES ($1)", clientIDBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to insert chain: %w", err)
		}

		if err = commonSQLX.Get(&row, chainIDSelectQuery+lockClause, clientIDBytes); err != nil {
			return nil, fmt.Errorf("failed to get chain id: %w", err)
		}
	}

	return row.ID, nil
}

func (sqlDB *SQLDB) GetUpdatesForType(dataType int, clientToken int64, fetchFolders bool, chainID int64, maxSize int) (bool, []SyncEntity, error) {
	var additionalCondition string
	if !fetchFolders {
		additionalCondition = "AND folder = false "
	}
	query := `SELECT * FROM entities
		WHERE chain_id = $1 AND data_type = $2 AND mtime > $3 ` + additionalCondition + `ORDER BY mtime LIMIT $4`

	entities := []SyncEntity{}
	if err := sqlDB.Select(&entities, query, chainID, dataType, clientToken, maxSize); err != nil {
		return false, nil, fmt.Errorf("failed to get entity updates: %w", err)
	}
	return len(entities) == maxSize, entities, nil
}

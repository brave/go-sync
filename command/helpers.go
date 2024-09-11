package command

import (
	"fmt"

	"github.com/brave/go-sync/cache"
	"github.com/brave/go-sync/datastore"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DBHelpers struct {
	dynamoDB             datastore.DynamoDatastore
	sqlDB                datastore.SQLDB
	Trx                  *sqlx.Tx
	clientID             string
	ChainID              int64
	variationHashDecimal float32
	ItemCounts           *ItemCounts
}

func NewDBHelpers(dynamoDB datastore.DynamoDatastore, sqlDB datastore.SQLDB, clientID string, cache *cache.Cache, initItemCounts bool) (*DBHelpers, error) {
	trx, err := sqlDB.Beginx()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}

	chainID, err := sqlDB.GetAndLockChainID(trx, clientID)
	if err != nil {
		trx.Rollback()
		return nil, err
	}
	variationHashDecimal := datastore.VariationHashDecimal(clientID)

	var itemCounts *ItemCounts
	if initItemCounts {
		itemCounts, err = getItemCounts(cache, dynamoDB, sqlDB, trx, clientID, *chainID)
		if err != nil {
			trx.Rollback()
			return nil, err
		}
	}

	return &DBHelpers{
		dynamoDB:             dynamoDB,
		sqlDB:                sqlDB,
		Trx:                  trx,
		clientID:             clientID,
		ChainID:              *chainID,
		variationHashDecimal: variationHashDecimal,
		ItemCounts:           itemCounts,
	}, nil
}

func (h *DBHelpers) hasItemInEitherDB(entity *datastore.SyncEntity) (exists bool, err error) {
	// Check if item exists using client_unique_tag
	if h.sqlDB.Variations.ShouldSaveToSQL(*entity.DataType, h.variationHashDecimal) {
		exists, err := h.sqlDB.HasItem(h.Trx, h.ChainID, *entity.ClientDefinedUniqueTag)
		if err != nil {
			return false, err
		}
		if !exists {
			return h.dynamoDB.HasItem(h.clientID, *entity.ClientDefinedUniqueTag)
		}
		return exists, err
	}
	return h.dynamoDB.HasItem(h.clientID, *entity.ClientDefinedUniqueTag)
}

func (h *DBHelpers) getUpdatesFromDBs(dataType int, token int64, fetchFolders bool, curMaxSize int) (hasChangesRemaining bool, syncEntities []datastore.SyncEntity, err error) {
	if curMaxSize == 0 {
		return false, nil, nil
	}
	if h.sqlDB.Variations.ShouldSaveToSQL(dataType, h.variationHashDecimal) {
		dynamoMigrationStatus, err := h.sqlDB.GetDynamoMigrationStatus(h.ChainID, dataType)
		if err != nil {
			return false, nil, err
		}

		if dynamoMigrationStatus == nil || dynamoMigrationStatus.EarliestMtime > token {
			var earliestMtime *int64
			if dynamoMigrationStatus != nil {
				earliestMtime = &dynamoMigrationStatus.EarliestMtime
			}
			hasChangesRemaining, syncEntities, err = h.dynamoDB.GetUpdatesForType(dataType, token, fetchFolders, h.clientID, curMaxSize, earliestMtime)
			if err != nil {
				return false, nil, err
			}
			curMaxSize -= len(syncEntities)
		}

		if curMaxSize > 0 {
			sqlHasChangesRemaining, sqlSyncEntities, err := h.sqlDB.GetUpdatesForType(dataType, token, fetchFolders, h.ChainID, curMaxSize)
			if err != nil {
				return false, nil, err
			}
			if sqlHasChangesRemaining {
				hasChangesRemaining = true
			}
			syncEntities = append(syncEntities, sqlSyncEntities...)
		}

		return hasChangesRemaining, syncEntities, nil
	}
	return h.dynamoDB.GetUpdatesForType(dataType, token, fetchFolders, h.clientID, curMaxSize, nil)
}

func (h *DBHelpers) insertSyncEntity(entity *datastore.SyncEntity) (conflict bool, err error) {
	savedInSQL := h.sqlDB.Variations.ShouldSaveToSQL(*entity.DataType, h.variationHashDecimal)
	if savedInSQL {
		conflict, err = h.sqlDB.InsertSyncEntity(h.Trx, entity)
	} else {
		conflict, err = h.dynamoDB.InsertSyncEntity(entity)
	}
	if err == nil && !conflict && (entity.Deleted == nil || !*entity.Deleted) {
		if err = h.ItemCounts.recordChange(*entity.DataType, false, savedInSQL); err != nil {
			return false, err
		}
	}
	return conflict, nil
}

func (h *DBHelpers) updateSyncEntity(entity *datastore.SyncEntity, oldVersion int64) (conflict bool, err error) {
	if h.sqlDB.Variations.ShouldSaveToSQL(*entity.DataType, h.variationHashDecimal) {
		conflict, err := h.sqlDB.UpdateSyncEntity(h.Trx, entity, oldVersion)
		if err != nil {
			return false, err
		}
		if conflict {
			oldEntity, err := h.dynamoDB.DeleteEntity(entity)
			if err != nil {
				return false, err
			}
			if oldEntity == nil {
				return true, nil
			}
			if oldEntity.Deleted == nil || !*oldEntity.Deleted {
				if err = h.ItemCounts.recordChange(*entity.DataType, true, false); err != nil {
					return false, err
				}
			}
			if *entity.DataType == datastore.HistoryTypeID {
				newID, err := uuid.NewV7()
				if err != nil {
					return false, err
				}
				entity.ID = newID.String()
			}
			conflict, err = h.sqlDB.InsertSyncEntity(h.Trx, entity)
			if err != nil {
				return false, err
			}
			if !conflict && (entity.Deleted == nil || !*entity.Deleted) {
				if err = h.ItemCounts.recordChange(*entity.DataType, false, true); err != nil {
					return false, err
				}
			}
			return conflict, err
		}
		return conflict, err
	}
	conflict, _, err = h.dynamoDB.UpdateSyncEntity(entity, oldVersion)
	return conflict, err
}

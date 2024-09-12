package command

import (
	"fmt"
	"math/rand/v2"

	"github.com/brave/go-sync/cache"
	"github.com/brave/go-sync/datastore"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DBHelpers struct {
	dynamoDB             datastore.DynamoDatastore
	sqlDB                datastore.SQLDatastore
	Trx                  *sqlx.Tx
	clientID             string
	ChainID              int64
	variationHashDecimal float32
	ItemCounts           *ItemCounts
}

func NewDBHelpers(dynamoDB datastore.DynamoDatastore, sqlDB datastore.SQLDatastore, clientID string, cache *cache.Cache, initItemCounts bool) (*DBHelpers, error) {
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
	if h.sqlDB.Variations().ShouldSaveToSQL(*entity.DataType, h.variationHashDecimal) {
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
	if h.sqlDB.Variations().ShouldSaveToSQL(dataType, h.variationHashDecimal) {
		dynamoMigrationStatuses, err := h.sqlDB.GetDynamoMigrationStatuses(h.Trx, h.ChainID, []int{dataType})
		if err != nil {
			return false, nil, err
		}

		if migrationStatus := dynamoMigrationStatuses[dataType]; migrationStatus == nil || (migrationStatus.EarliestMtime != nil && *migrationStatus.EarliestMtime > token) {
			var earliestMtime *int64
			if migrationStatus != nil {
				earliestMtime = migrationStatus.EarliestMtime
			}
			hasChangesRemaining, syncEntities, err = h.dynamoDB.GetUpdatesForType(dataType, &token, earliestMtime, fetchFolders, h.clientID, curMaxSize, true)
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
	return h.dynamoDB.GetUpdatesForType(dataType, &token, nil, fetchFolders, h.clientID, curMaxSize, true)
}

func (h *DBHelpers) insertSyncEntity(entity *datastore.SyncEntity) (conflict bool, err error) {
	savedInSQL := h.sqlDB.Variations().ShouldSaveToSQL(*entity.DataType, h.variationHashDecimal)
	if savedInSQL {
		conflict, err = h.sqlDB.InsertSyncEntities(h.Trx, []*datastore.SyncEntity{entity})
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
	if h.sqlDB.Variations().ShouldSaveToSQL(*entity.DataType, h.variationHashDecimal) {
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
			conflict, err = h.sqlDB.InsertSyncEntities(h.Trx, []*datastore.SyncEntity{entity})
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

func (h *DBHelpers) maybeMigrateToSQL(dataTypes []int) (migratedEntities []*datastore.SyncEntity, err error) {
	if rand.Float32() > h.sqlDB.MigrateIntervalPercent() {
		return nil, nil
	}
	var applicableDataTypes []int
	for _, dataType := range dataTypes {
		if !h.sqlDB.Variations().ShouldMigrateToSQL(dataType, h.variationHashDecimal) {
			continue
		}
		applicableDataTypes = append(applicableDataTypes, dataType)
	}
	migrationStatuses, err := h.sqlDB.GetDynamoMigrationStatuses(h.Trx, h.ChainID, applicableDataTypes)
	if err != nil {
		return nil, err
	}

	currLimit := h.sqlDB.MigrateChunkSize()
	var updatedMigrationStatuses []*datastore.MigrationStatus

	for _, dataType := range dataTypes {
		if currLimit <= 0 {
			break
		}
		migrationStatus := migrationStatuses[dataType]
		if migrationStatus != nil && migrationStatus.EarliestMtime == nil {
			continue
		}

		var earliestMtime *int64
		if migrationStatus != nil {
			earliestMtime = migrationStatus.EarliestMtime
		}

		hasChangesRemaining, syncEntities, err := h.dynamoDB.GetUpdatesForType(dataType, nil, earliestMtime, true, h.clientID, currLimit, false)
		if err != nil {
			return nil, err
		}

		currLimit -= len(syncEntities)

		lastItem := &syncEntities[len(syncEntities)-1]

		if !hasChangesRemaining {
			migrationStatus.EarliestMtime = nil
		} else if lastItem.Mtime != nil {
			migrationStatus.EarliestMtime = lastItem.Mtime
		}
		updatedMigrationStatuses = append(updatedMigrationStatuses, migrationStatus)

		var syncEntitiesPtr []*datastore.SyncEntity
		for _, syncEntity := range syncEntities {
			syncEntitiesPtr = append(syncEntitiesPtr, &syncEntity)
			migratedEntities = append(migratedEntities, &syncEntity)
		}

		if _, err = h.sqlDB.InsertSyncEntities(h.Trx, syncEntitiesPtr); err != nil {
			return nil, err
		}
	}
	if err = h.sqlDB.UpdateDynamoMigrationStatuses(h.Trx, updatedMigrationStatuses); err != nil {
		return nil, err
	}
	return migratedEntities, nil
}

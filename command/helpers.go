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
	SQLDB                datastore.SQLDatastore
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
		SQLDB:                sqlDB,
		Trx:                  trx,
		clientID:             clientID,
		ChainID:              *chainID,
		variationHashDecimal: variationHashDecimal,
		ItemCounts:           itemCounts,
	}, nil
}

func (h *DBHelpers) hasItemInEitherDB(entity *datastore.SyncEntity) (exists bool, err error) {
	// Check if item exists using client_unique_tag
	if h.SQLDB.Variations().ShouldSaveToSQL(*entity.DataType, h.variationHashDecimal) {
		exists, err := h.SQLDB.HasItem(h.Trx, h.ChainID, *entity.ClientDefinedUniqueTag)
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
	if h.SQLDB.Variations().ShouldSaveToSQL(dataType, h.variationHashDecimal) {
		dynamoMigrationStatuses, err := h.SQLDB.GetDynamoMigrationStatuses(h.Trx, h.ChainID, []int{dataType})
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
			sqlHasChangesRemaining, sqlSyncEntities, err := h.SQLDB.GetUpdatesForType(h.Trx, dataType, token, fetchFolders, h.ChainID, curMaxSize)
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
	savedInSQL := h.SQLDB.Variations().ShouldSaveToSQL(*entity.DataType, h.variationHashDecimal)
	if savedInSQL {
		conflict, err = h.SQLDB.InsertSyncEntities(h.Trx, []*datastore.SyncEntity{entity})
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

func getMigratedEntityID(entity *datastore.SyncEntity) (string, error) {
	id := entity.ID
	if *entity.DataType == datastore.HistoryTypeID {
		newID, err := uuid.NewV7()
		if err != nil {
			return "", err
		}
		id = newID.String()
	}
	return id, nil
}

func (h *DBHelpers) updateSyncEntity(entity *datastore.SyncEntity, oldVersion int64) (conflict bool, migratedEntity *datastore.SyncEntity, err error) {
	var deleted bool
	shouldSaveInSQL := h.SQLDB.Variations().ShouldSaveToSQL(*entity.DataType, h.variationHashDecimal)
	if shouldSaveInSQL {
		conflict, deleted, err = h.SQLDB.UpdateSyncEntity(h.Trx, entity, oldVersion)
		if err != nil {
			return false, nil, err
		}
		// Conflict might mean that the entity does not exist in SQL but exists in Dynamo.
		// Check for a Dynamo entity and migrate it accordingly.
		if conflict {
			oldEntity, err := h.dynamoDB.GetEntity(datastore.ItemQuery{
				ID:       entity.ID,
				ClientID: entity.ClientID,
			})
			if err != nil {
				return false, nil, err
			}
			if oldEntity == nil {
				// The conflict is unrelated to a pending Dynamo to SQL migration.
				// Return conflict error to client.
				return true, nil, nil
			}
			if oldEntity.Deleted == nil || !*oldEntity.Deleted {
				// If the stored entity was not already deleted, decrement the
				// Dynamo item count since we'll be migrating the entity to SQL.
				if err = h.ItemCounts.recordChange(*entity.DataType, true, false); err != nil {
					return false, nil, err
				}
			}
			migratedEntityID, err := getMigratedEntityID(entity)
			if err != nil {
				return false, nil, err
			}
			entity.ID = migratedEntityID
			conflict, err = h.SQLDB.InsertSyncEntities(h.Trx, []*datastore.SyncEntity{entity})
			if err != nil {
				return false, nil, err
			}
			if !conflict && (entity.Deleted == nil || !*entity.Deleted) {
				// If the new entity is not considered deleted, increment the
				// SQL interim count.
				if err = h.ItemCounts.recordChange(*entity.DataType, false, true); err != nil {
					return false, nil, err
				}
			}
			return conflict, oldEntity, err
		}
	} else {
		conflict, deleted, err = h.dynamoDB.UpdateSyncEntity(entity, oldVersion)
		if err != nil {
			return false, nil, err
		}
	}
	if !conflict && deleted {
		if err = h.ItemCounts.recordChange(*entity.DataType, true, shouldSaveInSQL); err != nil {
			return false, nil, err
		}
	}
	return conflict, nil, err
}

func (h *DBHelpers) maybeMigrateToSQL(dataTypes []int) (migratedEntities []*datastore.SyncEntity, err error) {
	if !h.SQLDB.Variations().Ready {
		return nil, nil
	}
	if rand.Float32() > h.SQLDB.MigrateIntervalPercent() {
		return nil, nil
	}
	var applicableDataTypes []int
	for _, dataType := range dataTypes {
		if !h.SQLDB.Variations().ShouldMigrateToSQL(dataType, h.variationHashDecimal) {
			continue
		}
		applicableDataTypes = append(applicableDataTypes, dataType)
	}
	if len(applicableDataTypes) == 0 {
		return nil, nil
	}

	migrationStatuses, err := h.SQLDB.GetDynamoMigrationStatuses(h.Trx, h.ChainID, applicableDataTypes)
	if err != nil {
		return nil, err
	}

	currLimit := h.SQLDB.MigrateChunkSize()
	var updatedMigrationStatuses []*datastore.MigrationStatus

	for _, dataType := range applicableDataTypes {
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
		} else {
			migrationStatus = &datastore.MigrationStatus{
				ChainID:       h.ChainID,
				DataType:      dataType,
				EarliestMtime: nil,
			}
		}

		hasChangesRemaining, syncEntities, err := h.dynamoDB.GetUpdatesForType(dataType, nil, earliestMtime, true, h.clientID, currLimit, false)
		if err != nil {
			return nil, err
		}

		currLimit -= len(syncEntities)

		if !hasChangesRemaining {
			migrationStatus.EarliestMtime = nil
		} else if len(syncEntities) > 0 {
			if lastItem := &syncEntities[len(syncEntities)-1]; lastItem.Mtime != nil {
				migrationStatus.EarliestMtime = lastItem.Mtime
			}
		}
		updatedMigrationStatuses = append(updatedMigrationStatuses, migrationStatus)

		var syncEntitiesPtr []*datastore.SyncEntity
		for _, syncEntity := range syncEntities {
			syncEntity.ChainID = &h.ChainID
			newEntity := &syncEntity
			migratedEntityID, err := getMigratedEntityID(&syncEntity)
			if err != nil {
				return nil, err
			}
			if migratedEntityID != syncEntity.ID {
				entityClone := syncEntity
				entityClone.ID = migratedEntityID
				newEntity = &entityClone
			}
			syncEntitiesPtr = append(syncEntitiesPtr, newEntity)
			migratedEntities = append(migratedEntities, &syncEntity)
		}

		if len(syncEntitiesPtr) > 0 {
			if _, err = h.SQLDB.InsertSyncEntities(h.Trx, syncEntitiesPtr); err != nil {
				return nil, err
			}
		}
	}
	if len(updatedMigrationStatuses) > 0 {
		if err = h.SQLDB.UpdateDynamoMigrationStatuses(h.Trx, updatedMigrationStatuses); err != nil {
			return nil, err
		}
	}
	return migratedEntities, nil
}

// InsertServerDefinedUniqueEntities inserts the server defined unique tag
// entities if it is not in the DB yet for a specific client.
func (h *DBHelpers) InsertServerDefinedUniqueEntities() error {
	if !h.SQLDB.Variations().Ready {
		return fmt.Errorf("SQL rollout not ready")
	}
	// Check if they're existed already for this client.
	// If yes, just return directly.
	ready, err := h.dynamoDB.HasServerDefinedUniqueTag(h.clientID, nigoriTag)
	if err != nil {
		return fmt.Errorf("error checking if entity with a server tag existed: %w", err)
	}
	if ready {
		return nil
	}

	entities, err := CreateServerDefinedUniqueEntities(h.clientID, h.ChainID)
	if err != nil {
		return err
	}

	var dynamoEntities []*datastore.SyncEntity
	var sqlEntities []*datastore.SyncEntity
	for _, entity := range entities {
		if h.SQLDB.Variations().ShouldSaveToSQL(*entity.DataType, h.variationHashDecimal) {
			sqlEntities = append(sqlEntities, entity)
		} else {
			dynamoEntities = append(dynamoEntities, entity)
		}
	}

	if len(dynamoEntities) > 0 {
		err = h.dynamoDB.InsertSyncEntitiesWithServerTags(dynamoEntities)
		if err != nil {
			return fmt.Errorf("error inserting entities with server tags to DynamoDB: %w", err)
		}
	}

	if len(sqlEntities) > 0 {
		_, err = h.SQLDB.InsertSyncEntities(h.Trx, sqlEntities)
		if err != nil {
			return fmt.Errorf("error inserting entities with server tags to SQL: %w", err)
		}
	}

	return nil
}

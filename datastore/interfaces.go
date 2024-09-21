package datastore

import "github.com/jmoiron/sqlx"

// DynamoDatastore abstracts over the underlying datastore.
type DynamoDatastore interface {
	// Insert a new sync entity.
	InsertSyncEntity(entity *SyncEntity) (bool, error)
	// Insert a series of sync entities in a write transaction.
	InsertSyncEntitiesWithServerTags(entities []*SyncEntity) error
	// Update an existing sync entity.
	UpdateSyncEntity(entity *SyncEntity, oldVersion int64) (conflict bool, deleted bool, err error)
	// Get updates for a specific type which are modified after the time of
	// client token for a given client. Besides the array of sync entities, a
	// boolean value indicating whether there are more updates to query in the
	// next batch is returned.
	GetUpdatesForType(dataType int, minMtime *int64, maxMtime *int64, fetchFolders bool, clientID string, maxSize int, ascOrder bool) (hasChangesRemaining bool, entities []SyncEntity, err error)
	// Check if a server-defined unique tag is in the datastore.
	HasServerDefinedUniqueTag(clientID string, tag string) (bool, error)
	// Get the count of sync items for a client.
	GetClientItemCount(clientID string) (*DynamoItemCounts, error)
	// Update the count of sync items for a client.
	UpdateClientItemCount(counts *DynamoItemCounts, newNormalItemCount int, newHistoryItemCount int) error
	// ClearServerData deletes all items for a given clientID
	ClearServerData(clientID string) ([]SyncEntity, error)
	// DisableSyncChain marks a chain as disabled so no further updates or commits can happen
	DisableSyncChain(clientID string) error
	// IsSyncChainDisabled checks whether a given sync chain is deleted
	IsSyncChainDisabled(clientID string) (bool, error)
	// HasItem checks if sync item exists for a client
	HasItem(clientID string, ID string) (bool, error)
	// GetEntity gets an existing entity
	GetEntity(query ItemQuery) (*SyncEntity, error)
	// DeleteEntities deletes multiple existing items
	DeleteEntities(entities []*SyncEntity) error
}

// SQLDatastore abstracts over the underlying datastore.
type SQLDatastore interface {
	// InsertSyncEntities inserts multiple sync entities into the database
	InsertSyncEntities(tx *sqlx.Tx, entities []*SyncEntity) (bool, error)
	// HasItem checks if an item exists in the database
	HasItem(tx *sqlx.Tx, chainID int64, clientTag string) (bool, error)
	// UpdateSyncEntity updates a sync entity in the database
	UpdateSyncEntity(tx *sqlx.Tx, entity *SyncEntity, oldVersion int64) (conflict bool, deleted bool, err error)
	// GetAndLockChainID retrieves and locks a chain ID for a given client ID
	GetAndLockChainID(tx *sqlx.Tx, clientID string) (*int64, error)
	// GetUpdatesForType retrieves updates for a specific data type
	GetUpdatesForType(tx *sqlx.Tx, dataType int, clientToken int64, fetchFolders bool, chainID int64, maxSize int) (hasChangesRemaining bool, entities []SyncEntity, err error)
	// GetDynamoMigrationStatuses retrieves migration statuses for specified data types
	GetDynamoMigrationStatuses(tx *sqlx.Tx, chainID int64, dataTypes []int) (map[int]*MigrationStatus, error)
	// UpdateDynamoMigrationStatuses updates migration statuses in the database
	UpdateDynamoMigrationStatuses(tx *sqlx.Tx, statuses []*MigrationStatus) error
	// GetItemCounts provides the counts of items associated with a chain
	GetItemCounts(tx *sqlx.Tx, chainID int64) (*SQLItemCounts, error)
	// Beginx initializes a database transaction
	Beginx() (*sqlx.Tx, error)
	// Variations returns the SQLVariations utility
	Variations() *SQLVariations
	// MigrateIntervalPercent returns the percentage of update requests that will perform
	// a chunked migration
	MigrateIntervalPercent() float32
	// MigrateChunkSize returns the max chunk size of migration attempts
	MigrateChunkSize() int
	// DeleteChain removes a chain and its associated data from the database
	DeleteChain(tx *sqlx.Tx, chainID int64) error
}

package datastore

// DynamoDatastore abstracts over the underlying datastore.
type DynamoDatastore interface {
	// Insert a new sync entity.
	InsertSyncEntity(entity *SyncEntity) (bool, error)
	// Insert a series of sync entities in a write transaction.
	InsertSyncEntitiesWithServerTags(entities []*SyncEntity) error
	// Update an existing sync entity.
	UpdateSyncEntity(entity *SyncEntity, oldVersion int64) (conflict bool, delete bool, err error)
	// Get updates for a specific type which are modified after the time of
	// client token for a given client. Besides the array of sync entities, a
	// boolean value indicating whether there are more updates to query in the
	// next batch is returned.
	GetUpdatesForType(dataType int, clientToken int64, fetchFolders bool, clientID string, maxSize int) (bool, []SyncEntity, error)
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
	// Checks if sync item exists for a client
	HasItem(clientID string, ID string) (bool, error)
}

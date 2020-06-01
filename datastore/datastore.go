package datastore

// Datastore abstracts over the underlying datastore.
type Datastore interface {
	// Insert a new sync entity.
	InsertSyncEntity(entity *SyncEntity) error
	// Insert a series of sync entities in a write transaction.
	InsertSyncEntitiesWithServerTags(entities []*SyncEntity) error
	// Update an existing sync entity.
	UpdateSyncEntity(entity *SyncEntity) (conflict bool, delete bool, err error)
	// Get updates for a specific type which are modified after the time of
	// client token for a given client. The total count of the updates available
	// in the DB would also be returned, which might not be the same as the
	// length of the returned SyncEntity slice if the total count is greater than
	// the maxSize parameter, in that case, length of the returned SyncEntity
	// slice would be maxSize.
	GetUpdatesForType(dataType int, clientToken int64, fetchFolders bool, clientID string, maxSize int64) (int64, []SyncEntity, error)
	// Check if a server-defined unique tag is in the datastore.
	HasServerDefinedUniqueTag(clientID string, tag string) (bool, error)
	// Get the count of sync items for a client.
	GetClientItemCount(clientID string) (int, error)
	// Update the count of sync items for a client.
	UpdateClientItemCount(clientID string, count int) error
}

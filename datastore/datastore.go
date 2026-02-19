package datastore

import "context"

// Datastore abstracts over the underlying datastore.
type Datastore interface {
	// Insert a new sync entity.
	InsertSyncEntity(ctx context.Context, entity *SyncEntity) (bool, error)
	// Insert a series of sync entities in a write transaction.
	InsertSyncEntitiesWithServerTags(ctx context.Context, entities []*SyncEntity) error
	// Update an existing sync entity.
	UpdateSyncEntity(ctx context.Context, entity *SyncEntity, oldVersion int64) (conflict bool, deleted bool, err error)
	// Get updates for a specific type which are modified after the time of
	// client token for a given client. Besides the array of sync entities, a
	// boolean value indicating whether there are more updates to query in the
	// next batch is returned.
	GetUpdatesForType(ctx context.Context, dataType int, clientToken int64, fetchFolders bool, clientID string, maxSize int64) (bool, []SyncEntity, error)
	// Check if a server-defined unique tag is in the datastore.
	HasServerDefinedUniqueTag(ctx context.Context, clientID string, tag string) (bool, error)
	// Get the count of sync items for a client.
	GetClientItemCount(ctx context.Context, clientID string) (*ClientItemCounts, error)
	// Update the count of sync items for a client.
	UpdateClientItemCount(ctx context.Context, counts *ClientItemCounts, newNormalItemCount int, newHistoryItemCount int) error
	// ClearServerData deletes all items for a given clientID
	ClearServerData(ctx context.Context, clientID string) ([]SyncEntity, error)
	// DisableSyncChain marks a chain as disabled so no further updates or commits can happen
	DisableSyncChain(ctx context.Context, clientID string) error
	// IsSyncChainDisabled checks whether a given sync chain is deleted
	IsSyncChainDisabled(ctx context.Context, clientID string) (bool, error)
	// Checks if sync item exists for a client
	HasItem(ctx context.Context, clientID string, ID string) (bool, error)
}

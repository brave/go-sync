package datastoretest

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/brave/go-sync/datastore"
)

// MockDatastore is used to mock datastorein tests
type MockDatastore struct {
	mock.Mock
}

// InsertSyncEntity mocks calls to InsertSyncEntity
func (m *MockDatastore) InsertSyncEntity(ctx context.Context, entity *datastore.SyncEntity) (bool, error) {
	args := m.Called(ctx, entity)
	return args.Bool(0), args.Error(1)
}

// InsertSyncEntitiesWithServerTags mocks calls to InsertSyncEntitiesWithServerTags
func (m *MockDatastore) InsertSyncEntitiesWithServerTags(ctx context.Context, entities []*datastore.SyncEntity) error {
	args := m.Called(ctx, entities)
	return args.Error(0)
}

// UpdateSyncEntity mocks calls to UpdateSyncEntity
func (m *MockDatastore) UpdateSyncEntity(ctx context.Context, entity *datastore.SyncEntity, oldVersion int64) (conflict bool, deleted bool, err error) {
	args := m.Called(ctx, entity, oldVersion)
	return args.Bool(0), args.Bool(1), args.Error(2)
}

// GetUpdatesForType mocks calls to GetUpdatesForType
func (m *MockDatastore) GetUpdatesForType(ctx context.Context, dataType int, clientToken int64, fetchFolders bool, clientID string, maxSize int64) (bool, []datastore.SyncEntity, error) {
	args := m.Called(ctx, dataType, clientToken, fetchFolders, clientID, maxSize)
	return args.Bool(0), args.Get(1).([]datastore.SyncEntity), args.Error(2)
}

// HasServerDefinedUniqueTag mocks calls to HasServerDefinedUniqueTag
func (m *MockDatastore) HasServerDefinedUniqueTag(ctx context.Context, clientID string, tag string) (bool, error) {
	args := m.Called(ctx, clientID, tag)
	return args.Bool(0), args.Error(1)
}

func (m *MockDatastore) HasItem(ctx context.Context, clientID string, ID string) (bool, error) {
	args := m.Called(ctx, clientID, ID)
	return args.Bool(0), args.Error(1)
}

// GetClientItemCount mocks calls to GetClientItemCount
func (m *MockDatastore) GetClientItemCount(ctx context.Context, clientID string) (*datastore.ClientItemCounts, error) {
	args := m.Called(ctx, clientID)
	return &datastore.ClientItemCounts{ClientID: clientID, ID: clientID}, args.Error(1)
}

// UpdateClientItemCount mocks calls to UpdateClientItemCount
func (m *MockDatastore) UpdateClientItemCount(ctx context.Context, counts *datastore.ClientItemCounts, newNormalItemCount int, newHistoryItemCount int) error {
	args := m.Called(ctx, counts, newNormalItemCount, newHistoryItemCount)
	return args.Error(0)
}

// ClearServerData mocks calls to ClearServerData
func (m *MockDatastore) ClearServerData(ctx context.Context, clientID string) ([]datastore.SyncEntity, error) {
	args := m.Called(ctx, clientID)
	return args.Get(0).([]datastore.SyncEntity), args.Error(1)
}

// DisableSyncChain mocks calls to DisableSyncChain
func (m *MockDatastore) DisableSyncChain(ctx context.Context, clientID string) error {
	args := m.Called(ctx, clientID)
	return args.Error(0)
}

// IsSyncChainDisabled mocks calls to IsSyncChainDisabled
func (m *MockDatastore) IsSyncChainDisabled(ctx context.Context, clientID string) (bool, error) {
	args := m.Called(ctx, clientID)
	return args.Bool(0), args.Error(1)
}

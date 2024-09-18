package datastoretest

import (
	"github.com/brave/go-sync/datastore"
	"github.com/stretchr/testify/mock"
)

// MockDatastore is used to mock datastorein tests
type MockDatastore struct {
	mock.Mock
}

// InsertSyncEntity mocks calls to InsertSyncEntity
func (m *MockDatastore) InsertSyncEntity(entity *datastore.SyncEntity) (bool, error) {
	args := m.Called(entity)
	return args.Bool(0), args.Error(1)
}

// InsertSyncEntitiesWithServerTags mocks calls to InsertSyncEntitiesWithServerTags
func (m *MockDatastore) InsertSyncEntitiesWithServerTags(entities []*datastore.SyncEntity) error {
	args := m.Called(entities)
	return args.Error(0)
}

// UpdateSyncEntity mocks calls to UpdateSyncEntity
func (m *MockDatastore) UpdateSyncEntity(entity *datastore.SyncEntity, oldVersion int64) (conflict bool, err error) {
	args := m.Called(entity, oldVersion)
	return args.Bool(0), args.Error(1)
}

// GetUpdatesForType mocks calls to GetUpdatesForType
func (m *MockDatastore) GetUpdatesForType(dataType int, minMtime *int64, maxMtime *int64, fetchFolders bool, clientID string, maxSize int, ascOrder bool) (bool, []datastore.SyncEntity, error) {
	args := m.Called(dataType, minMtime, maxMtime, fetchFolders, clientID, maxSize, ascOrder)
	return args.Bool(0), args.Get(1).([]datastore.SyncEntity), args.Error(2)
}

// HasServerDefinedUniqueTag mocks calls to HasServerDefinedUniqueTag
func (m *MockDatastore) HasServerDefinedUniqueTag(clientID string, tag string) (bool, error) {
	args := m.Called(clientID, tag)
	return args.Bool(0), args.Error(1)
}

func (m *MockDatastore) HasItem(clientID string, ID string) (bool, error) {
	args := m.Called(clientID, ID)
	return args.Bool(0), args.Error(1)
}

// GetClientItemCount mocks calls to GetClientItemCount
func (m *MockDatastore) GetClientItemCount(clientID string) (*datastore.DynamoItemCounts, error) {
	args := m.Called(clientID)
	return &datastore.DynamoItemCounts{ClientID: clientID, ID: clientID}, args.Error(1)
}

// UpdateClientItemCount mocks calls to UpdateClientItemCount
func (m *MockDatastore) UpdateClientItemCount(counts *datastore.DynamoItemCounts, newNormalItemCount int, newHistoryItemCount int) error {
	args := m.Called(counts, newNormalItemCount, newHistoryItemCount)
	return args.Error(0)
}

// ClearServerData mocks calls to ClearServerData
func (m *MockDatastore) ClearServerData(clientID string) ([]datastore.SyncEntity, error) {
	args := m.Called(clientID)
	return args.Get(0).([]datastore.SyncEntity), args.Error(1)
}

// DisableSyncChain mocks calls to DisableSyncChain
func (m *MockDatastore) DisableSyncChain(clientID string) error {
	args := m.Called(clientID)
	return args.Error(0)
}

// IsSyncChainDisabled mocks calls to IsSyncChainDisabled
func (m *MockDatastore) IsSyncChainDisabled(clientID string) (bool, error) {
	args := m.Called(clientID)
	return args.Bool(0), args.Error(1)
}

// DeleteEntities mocks the deletion of sync entities
func (m *MockDatastore) DeleteEntities(entities []*datastore.SyncEntity) error {
	args := m.Called(entities)
	return args.Error(0)
}

// GetEntity mocks the retrieval of a sync entity
func (m *MockDatastore) GetEntity(query datastore.ItemQuery) (*datastore.SyncEntity, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*datastore.SyncEntity), args.Error(1)
}

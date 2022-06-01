package datastoretest

import (
	"github.com/brave/go-sync/datastore"
	"github.com/stretchr/testify/mock"
)

type MockDatastore struct {
	mock.Mock
}

func (m *MockDatastore) InsertSyncEntity(entity *datastore.SyncEntity) (bool, error) {
	args := m.Called(entity)
	return args.Bool(0), args.Error(1)
}

func (m *MockDatastore) InsertSyncEntitiesWithServerTags(entities []*datastore.SyncEntity) error {
	args := m.Called(entities)
	return args.Error(0)
}

func (m *MockDatastore) UpdateSyncEntity(entity *datastore.SyncEntity, oldVersion int64) (conflict bool, delete bool, err error) {
	args := m.Called(entity, oldVersion)
	return args.Bool(0), args.Bool(1), args.Error(2)
}

func (m *MockDatastore) GetUpdatesForType(dataType int, clientToken int64, fetchFolders bool, clientID string, maxSize int64) (bool, []datastore.SyncEntity, error) {
	args := m.Called(dataType, clientToken, fetchFolders, clientID, maxSize)
	return args.Bool(0), args.Get(1).([]datastore.SyncEntity), args.Error(2)
}

func (m *MockDatastore) HasServerDefinedUniqueTag(clientID string, tag string) (bool, error) {
	args := m.Called(clientID, tag)
	return args.Bool(0), args.Error(1)
}

func (m *MockDatastore) GetClientItemCount(clientID string) (int, error) {
	args := m.Called(clientID)
	return args.Int(0), args.Error(1)
}

func (m *MockDatastore) UpdateClientItemCount(clientID string, count int) error {
	args := m.Called(clientID, count)
	return args.Error(0)
}

func (m *MockDatastore) DeleteClientItems(clientID string) error {
	args := m.Called(clientID)
	return args.Error(0)
}

func (m *MockDatastore) IsSyncChainDisabled(clientID string) (bool, error) {
	args := m.Called(clientID)
	return args.Bool(0), args.Error(1)
}

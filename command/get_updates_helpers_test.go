package command //nolint:testpackage // tests unexported GetUpdates helpers

import (
	"context"
	"encoding/binary"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/brave/go-sync/cache"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/datastore/datastoretest"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
)

func TestIsNewClientAndPollOrigin(t *testing.T) {
	newClient := sync_pb.SyncEnums_NEW_CLIENT
	poll := sync_pb.SyncEnums_PERIODIC
	trigger := sync_pb.SyncEnums_GU_TRIGGER

	assert.False(t, isNewClientOrigin(&sync_pb.GetUpdatesMessage{}))
	assert.False(t, isPollOrigin(&sync_pb.GetUpdatesMessage{}))
	assert.True(t, isNewClientOrigin(&sync_pb.GetUpdatesMessage{GetUpdatesOrigin: &newClient}))
	assert.True(t, isPollOrigin(&sync_pb.GetUpdatesMessage{GetUpdatesOrigin: &poll}))
	assert.False(t, isNewClientOrigin(&sync_pb.GetUpdatesMessage{GetUpdatesOrigin: &trigger}))
	assert.False(t, isPollOrigin(&sync_pb.GetUpdatesMessage{GetUpdatesOrigin: &trigger}))
}

func TestShouldFetchFolders(t *testing.T) {
	assert.True(t, shouldFetchFolders(&sync_pb.GetUpdatesMessage{}))
	assert.True(t, shouldFetchFolders(&sync_pb.GetUpdatesMessage{FetchFolders: aws.Bool(true)}))
	assert.False(t, shouldFetchFolders(&sync_pb.GetUpdatesMessage{FetchFolders: aws.Bool(false)}))
}

func TestProgressMarkerToken(t *testing.T) {
	existing := []byte{1, 2, 3}
	assert.Equal(t, existing, progressMarkerToken(&sync_pb.DataTypeProgressMarker{Token: existing}))

	decoded, n := binary.Varint(progressMarkerToken(&sync_pb.DataTypeProgressMarker{}))
	require.Positive(t, n)
	assert.Equal(t, int64(0), decoded)
}

func TestMaybeSetNigoriEncryptionKeys(t *testing.T) {
	rsp := &sync_pb.GetUpdatesResponse{}
	maybeSetNigoriEncryptionKeys(rsp, nigoriTypeID, false)
	assert.Nil(t, rsp.EncryptionKeys)

	maybeSetNigoriEncryptionKeys(rsp, 32904, true)
	assert.Nil(t, rsp.EncryptionKeys)

	maybeSetNigoriEncryptionKeys(rsp, nigoriTypeID, true)
	require.Len(t, rsp.EncryptionKeys, 1)
	assert.Equal(t, []byte("1234"), rsp.EncryptionKeys[0])
}

func TestAppendGetUpdatesEntities(t *testing.T) {
	mtime := int64(42)
	entities := []datastore.SyncEntity{{
		ID:      "id1",
		Version: aws.Int64(1),
		Mtime:   &mtime,
		Ctime:   aws.Int64(1),
		Deleted: aws.Bool(false),
		Folder:  aws.Bool(false),
	}}
	rsp := &sync_pb.GetUpdatesResponse{
		Entries:           make([]*sync_pb.SyncEntity, 0, 10),
		NewProgressMarker: []*sync_pb.DataTypeProgressMarker{{}},
	}

	j, errCode, err := appendGetUpdatesEntities(rsp, 0, entities)
	require.NoError(t, err)
	assert.Nil(t, errCode)
	assert.Equal(t, 1, j)
	require.Len(t, rsp.Entries, 1)
	decoded, n := binary.Varint(rsp.NewProgressMarker[0].Token)
	require.Positive(t, n)
	assert.Equal(t, mtime, decoded)

	j, errCode, err = appendGetUpdatesEntities(rsp, 0, nil)
	require.NoError(t, err)
	assert.Nil(t, errCode)
	assert.Equal(t, 0, j)
}

func TestMaybeSetupNewClient_Skip(t *testing.T) {
	errCode, err := maybeSetupNewClient(context.Background(), new(datastoretest.MockDatastore), "client", false)
	require.NoError(t, err)
	assert.Nil(t, errCode)
}

func TestMaybeSetTypeMtime(t *testing.T) {
	c := cache.NewCache(cache.NewRedisClient())
	t.Cleanup(func() {
		require.NoError(t, c.FlushAll(context.Background()))
	})

	mtime := int64(99)
	entities := []datastore.SyncEntity{{Mtime: &mtime}}
	key := c.GetTypeMtimeKey("client", 1)

	maybeSetTypeMtime(context.Background(), c, "client", 1, 1, 1, 5, entities)
	val, err := c.Get(context.Background(), key, false)
	require.NoError(t, err)
	assert.Empty(t, val, "cache should not be updated when changes remaining = 1")

	maybeSetTypeMtime(context.Background(), c, "client", 1, 0, 0, 5, entities)
	val, err = c.Get(context.Background(), key, false)
	require.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(5, 10), val)

	maybeSetTypeMtime(context.Background(), c, "client", 1, 0, 1, 5, entities)
	val, err = c.Get(context.Background(), key, false)
	require.NoError(t, err)
	assert.Equal(t, strconv.FormatInt(mtime, 10), val)
}

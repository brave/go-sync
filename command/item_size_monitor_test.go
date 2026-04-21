package command_test

import (
	"strings"
	"testing"

	"github.com/brave/go-sync/command"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func makeEntity(size int) *sync_pb.SyncEntity {
	return &sync_pb.SyncEntity{
		Name: proto.String(strings.Repeat("a", size)),
	}
}

func TestItemSizeMonitor_BelowThreshold(t *testing.T) {
	m := command.NewItemSizeMonitor()
	m.Observe(1, makeEntity(100))
	assert.Empty(t, m.StatsByType)
}

func TestItemSizeMonitor_AboveThreshold(t *testing.T) {
	m := command.NewItemSizeMonitor()
	m.Observe(1, makeEntity(command.LargeSizeThreshold+1))
	assert.Len(t, m.StatsByType, 1)
	assert.Equal(t, 1, m.StatsByType[1].Count)
}

func TestItemSizeMonitor_TracksMaxSize(t *testing.T) {
	m := command.NewItemSizeMonitor()
	small := makeEntity(command.LargeSizeThreshold + 1)
	large := makeEntity(command.LargeSizeThreshold + 10000)
	m.Observe(1, small)
	m.Observe(1, large)
	assert.Equal(t, 2, m.StatsByType[1].Count)
	assert.Equal(t, proto.Size(large), m.StatsByType[1].MaxSize)
}

func TestItemSizeMonitor_TracksPerDataType(t *testing.T) {
	m := command.NewItemSizeMonitor()
	entry := makeEntity(command.LargeSizeThreshold + 1)
	m.Observe(1, entry)
	m.Observe(2, entry)
	m.Observe(2, entry)
	assert.Equal(t, 1, m.StatsByType[1].Count)
	assert.Equal(t, 2, m.StatsByType[2].Count)
}

func TestItemSizeMonitor_LogWarningsOnlyForNonZero(t *testing.T) {
	m := command.NewItemSizeMonitor()
	m.Observe(1, makeEntity(100))
	m.Observe(2, makeEntity(command.LargeSizeThreshold+1))
	assert.NotContains(t, m.StatsByType, 1)
	assert.Contains(t, m.StatsByType, 2)
}

package command

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

const LargeSizeThreshold = 400 * 1024 // 400KB

type DataTypeStats struct {
	Count   int
	MaxSize int
}

// ItemSizeMonitor tracks, per data type, the count and max size of commit
// entries that exceed LargeSizeThreshold.
type ItemSizeMonitor struct {
	StatsByType map[int]*DataTypeStats
}

func NewItemSizeMonitor() *ItemSizeMonitor {
	return &ItemSizeMonitor{StatsByType: make(map[int]*DataTypeStats)}
}

// Observe measures the serialized size of msg and records it against dataType
// if it exceeds the threshold.
func (m *ItemSizeMonitor) Observe(dataType int, msg proto.Message) {
	size := proto.Size(msg)
	if size <= LargeSizeThreshold {
		return
	}
	stats := m.StatsByType[dataType]
	if stats == nil {
		stats = &DataTypeStats{}
		m.StatsByType[dataType] = stats
	}
	stats.Count++
	if size > stats.MaxSize {
		stats.MaxSize = size
	}
}

// LogWarnings emits a warning log for each data type with at least one entry
// that exceeded the threshold.
func (m *ItemSizeMonitor) LogWarnings() {
	for dataType, stats := range m.StatsByType {
		log.Warn().
			Int("data_type", dataType).
			Int("count", stats.Count).
			Int("max_size", stats.MaxSize).
			Int("threshold_bytes", LargeSizeThreshold).
			Msg("Commit entries exceeded size threshold")
	}
}

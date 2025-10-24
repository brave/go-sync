package datastore

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../.prom-gowrap.tmpl template

//go:generate gowrap gen -p github.com/brave/go-sync/datastore -i Datastore -t ../.prom-gowrap.tmpl -o instrumented_datastore.go

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// DatastoreWithPrometheus implements Datastore interface with all methods wrapped
// with Prometheus metrics
type DatastoreWithPrometheus struct {
	base         Datastore
	instanceName string
}

var datastoreDurationSummaryVec = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "datastore_duration_seconds",
		Help:       "datastore runtime duration and result",
		MaxAge:     time.Minute,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"instance_name", "method", "result"})

// NewDatastoreWithPrometheus returns an instance of the Datastore decorated with prometheus summary metric
func NewDatastoreWithPrometheus(base Datastore, instanceName string) DatastoreWithPrometheus {
	return DatastoreWithPrometheus{
		base:         base,
		instanceName: instanceName,
	}
}

// ClearServerData implements Datastore
func (_d DatastoreWithPrometheus) ClearServerData(ctx context.Context, clientID string) (sa1 []SyncEntity, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "ClearServerData", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.ClearServerData(ctx, clientID)
}

// DisableSyncChain implements Datastore
func (_d DatastoreWithPrometheus) DisableSyncChain(ctx context.Context, clientID string) (err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "DisableSyncChain", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.DisableSyncChain(ctx, clientID)
}

// GetClientItemCount implements Datastore
func (_d DatastoreWithPrometheus) GetClientItemCount(ctx context.Context, clientID string) (counts *ClientItemCounts, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "GetClientItemCount", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.GetClientItemCount(ctx, clientID)
}

// GetUpdatesForType implements Datastore
func (_d DatastoreWithPrometheus) GetUpdatesForType(ctx context.Context, dataType int, clientToken int64, fetchFolders bool, clientID string, maxSize int64) (b1 bool, sa1 []SyncEntity, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "GetUpdatesForType", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.GetUpdatesForType(ctx, dataType, clientToken, fetchFolders, clientID, maxSize)
}

// HasItem implements Datastore
func (_d DatastoreWithPrometheus) HasItem(ctx context.Context, clientID string, ID string) (b1 bool, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "HasItem", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.HasItem(ctx, clientID, ID)
}

// HasServerDefinedUniqueTag implements Datastore
func (_d DatastoreWithPrometheus) HasServerDefinedUniqueTag(ctx context.Context, clientID string, tag string) (b1 bool, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "HasServerDefinedUniqueTag", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.HasServerDefinedUniqueTag(ctx, clientID, tag)
}

// InsertSyncEntitiesWithServerTags implements Datastore
func (_d DatastoreWithPrometheus) InsertSyncEntitiesWithServerTags(ctx context.Context, entities []*SyncEntity) (err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "InsertSyncEntitiesWithServerTags", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.InsertSyncEntitiesWithServerTags(ctx, entities)
}

// InsertSyncEntity implements Datastore
func (_d DatastoreWithPrometheus) InsertSyncEntity(ctx context.Context, entity *SyncEntity) (b1 bool, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "InsertSyncEntity", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.InsertSyncEntity(ctx, entity)
}

// IsSyncChainDisabled implements Datastore
func (_d DatastoreWithPrometheus) IsSyncChainDisabled(ctx context.Context, clientID string) (b1 bool, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "IsSyncChainDisabled", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.IsSyncChainDisabled(ctx, clientID)
}

// UpdateClientItemCount implements Datastore
func (_d DatastoreWithPrometheus) UpdateClientItemCount(ctx context.Context, counts *ClientItemCounts, newNormalItemCount int, newHistoryItemCount int) (err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "UpdateClientItemCount", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.UpdateClientItemCount(ctx, counts, newNormalItemCount, newHistoryItemCount)
}

// UpdateSyncEntity implements Datastore
func (_d DatastoreWithPrometheus) UpdateSyncEntity(ctx context.Context, entity *SyncEntity, oldVersion int64) (conflict bool, deleted bool, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "UpdateSyncEntity", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.UpdateSyncEntity(ctx, entity, oldVersion)
}

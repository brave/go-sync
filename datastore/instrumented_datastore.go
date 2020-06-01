package datastore

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../.prom-gowrap.tmpl template

//go:generate gowrap gen -p github.com/brave/go-sync/datastore -i Datastore -t ../.prom-gowrap.tmpl -o instrumented_datastore.go

import (
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

// GetClientItemCount implements Datastore
func (_d DatastoreWithPrometheus) GetClientItemCount(clientID string) (i1 int, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "GetClientItemCount", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.GetClientItemCount(clientID)
}

// GetUpdatesForType implements Datastore
func (_d DatastoreWithPrometheus) GetUpdatesForType(dataType int, clientToken int64, fetchFolders bool, clientID string, maxSize int64) (i1 int64, sa1 []SyncEntity, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "GetUpdatesForType", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.GetUpdatesForType(dataType, clientToken, fetchFolders, clientID, maxSize)
}

// HasServerDefinedUniqueTag implements Datastore
func (_d DatastoreWithPrometheus) HasServerDefinedUniqueTag(clientID string, tag string) (b1 bool, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "HasServerDefinedUniqueTag", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.HasServerDefinedUniqueTag(clientID, tag)
}

// InsertSyncEntitiesWithServerTags implements Datastore
func (_d DatastoreWithPrometheus) InsertSyncEntitiesWithServerTags(entities []*SyncEntity) (err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "InsertSyncEntitiesWithServerTags", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.InsertSyncEntitiesWithServerTags(entities)
}

// InsertSyncEntity implements Datastore
func (_d DatastoreWithPrometheus) InsertSyncEntity(entity *SyncEntity) (err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "InsertSyncEntity", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.InsertSyncEntity(entity)
}

// UpdateClientItemCount implements Datastore
func (_d DatastoreWithPrometheus) UpdateClientItemCount(clientID string, count int) (err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "UpdateClientItemCount", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.UpdateClientItemCount(clientID, count)
}

// UpdateSyncEntity implements Datastore
func (_d DatastoreWithPrometheus) UpdateSyncEntity(entity *SyncEntity) (conflict bool, delete bool, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		datastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "UpdateSyncEntity", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.UpdateSyncEntity(entity)
}

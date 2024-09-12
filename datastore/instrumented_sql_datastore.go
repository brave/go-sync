// Code generated by gowrap. DO NOT EDIT.
// template: ../.prom-gowrap.tmpl
// gowrap: http://github.com/hexdigest/gowrap

package datastore

//go:generate gowrap gen -p github.com/brave/go-sync/datastore -i SQLDatastore -t ../.prom-gowrap.tmpl -o instrumented_sql_datastore.go -l ""

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// SQLDatastoreWithPrometheus implements SQLDatastore interface with all methods wrapped
// with Prometheus metrics
type SQLDatastoreWithPrometheus struct {
	base         SQLDatastore
	instanceName string
}

var sqldatastoreDurationSummaryVec = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "sqldatastore_duration_seconds",
		Help:       "sqldatastore runtime duration and result",
		MaxAge:     time.Minute,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"instance_name", "method", "result"})

// NewSQLDatastoreWithPrometheus returns an instance of the SQLDatastore decorated with prometheus summary metric
func NewSQLDatastoreWithPrometheus(base SQLDatastore, instanceName string) SQLDatastoreWithPrometheus {
	return SQLDatastoreWithPrometheus{
		base:         base,
		instanceName: instanceName,
	}
}

// Beginx implements SQLDatastore
func (_d SQLDatastoreWithPrometheus) Beginx() (tp1 *sqlx.Tx, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		sqldatastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "Beginx", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.Beginx()
}

// GetAndLockChainID implements SQLDatastore
func (_d SQLDatastoreWithPrometheus) GetAndLockChainID(tx *sqlx.Tx, clientID string) (ip1 *int64, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		sqldatastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "GetAndLockChainID", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.GetAndLockChainID(tx, clientID)
}

// GetDynamoMigrationStatuses implements SQLDatastore
func (_d SQLDatastoreWithPrometheus) GetDynamoMigrationStatuses(tx *sqlx.Tx, chainID int64, dataTypes []int) (m1 map[int]*MigrationStatus, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		sqldatastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "GetDynamoMigrationStatuses", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.GetDynamoMigrationStatuses(tx, chainID, dataTypes)
}

// GetItemCounts implements SQLDatastore
func (_d SQLDatastoreWithPrometheus) GetItemCounts(tx *sqlx.Tx, chainID int64) (sp1 *SQLItemCounts, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		sqldatastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "GetItemCounts", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.GetItemCounts(tx, chainID)
}

// GetUpdatesForType implements SQLDatastore
func (_d SQLDatastoreWithPrometheus) GetUpdatesForType(dataType int, clientToken int64, fetchFolders bool, chainID int64, maxSize int) (b1 bool, sa1 []SyncEntity, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		sqldatastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "GetUpdatesForType", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.GetUpdatesForType(dataType, clientToken, fetchFolders, chainID, maxSize)
}

// HasItem implements SQLDatastore
func (_d SQLDatastoreWithPrometheus) HasItem(tx *sqlx.Tx, chainId int64, clientTag string) (b1 bool, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		sqldatastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "HasItem", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.HasItem(tx, chainId, clientTag)
}

// InsertSyncEntities implements SQLDatastore
func (_d SQLDatastoreWithPrometheus) InsertSyncEntities(tx *sqlx.Tx, entities []*SyncEntity) (b1 bool, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		sqldatastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "InsertSyncEntities", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.InsertSyncEntities(tx, entities)
}

// MigrateChunkSize implements SQLDatastore
func (_d SQLDatastoreWithPrometheus) MigrateChunkSize() (i1 int) {
	_since := time.Now()
	defer func() {
		result := "ok"
		sqldatastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "MigrateChunkSize", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.MigrateChunkSize()
}

// MigrateIntervalPercent implements SQLDatastore
func (_d SQLDatastoreWithPrometheus) MigrateIntervalPercent() (f1 float32) {
	_since := time.Now()
	defer func() {
		result := "ok"
		sqldatastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "MigrateIntervalPercent", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.MigrateIntervalPercent()
}

// UpdateDynamoMigrationStatuses implements SQLDatastore
func (_d SQLDatastoreWithPrometheus) UpdateDynamoMigrationStatuses(tx *sqlx.Tx, statuses []*MigrationStatus) (err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		sqldatastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "UpdateDynamoMigrationStatuses", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.UpdateDynamoMigrationStatuses(tx, statuses)
}

// UpdateSyncEntity implements SQLDatastore
func (_d SQLDatastoreWithPrometheus) UpdateSyncEntity(tx *sqlx.Tx, entity *SyncEntity, oldVersion int64) (b1 bool, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		sqldatastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "UpdateSyncEntity", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.UpdateSyncEntity(tx, entity, oldVersion)
}

// Variations implements SQLDatastore
func (_d SQLDatastoreWithPrometheus) Variations() (sp1 *SQLVariations) {
	_since := time.Now()
	defer func() {
		result := "ok"
		sqldatastoreDurationSummaryVec.WithLabelValues(_d.instanceName, "Variations", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.Variations()
}

package datastore

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type SQLItemCounts struct {
	NormalItemCount  int `db:"normal_item_count"`
	HistoryItemCount int `db:"history_item_count"`
}

func (sqlDB *SQLDB) GetItemCounts(tx *sqlx.Tx, chainID int64) (*SQLItemCounts, error) {
	counts := SQLItemCounts{}
	err := tx.Get(&counts, `
		SELECT
			COUNT(*) FILTER (WHERE NOT (data_type = ANY($1))) AS normal_item_count,
			COUNT(*) FILTER (WHERE data_type = ANY($1)) AS history_item_count
		FROM entities
		WHERE chain_id = $2 AND deleted = false
	`, pq.Array([]int{HistoryTypeID, HistoryDeleteDirectiveTypeID}), chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get item counts: %w", err)
	}
	return &counts, nil
}

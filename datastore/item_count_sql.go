package datastore

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SQLItemCounts struct {
	NormalItemCount  int `db:"normal_item_count"`
	HistoryItemCount int `db:"history_item_count"`
}

func (sqlDB *SQLDB) GetItemCounts(tx *sqlx.Tx, chainID int64) (*SQLItemCounts, error) {
	counts := SQLItemCounts{}
	err := tx.Get(&counts, `
		SELECT
			COUNT(*) FILTER (WHERE data_type != $1) AS normal_item_count,
			COUNT(*) FILTER (WHERE data_type = $1) AS history_item_count
		FROM entities
		WHERE chain_id = $2 AND deleted = false
	`, HistoryTypeID, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get item counts: %w", err)
	}
	return &counts, nil
}

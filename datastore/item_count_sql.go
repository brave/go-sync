package datastore

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SQLItemCounts struct {
	NormalItemCount  int `db:"normal_item_count"`
	HistoryItemCount int `db:"history_item_count"`
}

// GetItemCounts returns the counts of items in the SQL database for a given chain ID
func (sqlDB *SQLDB) GetItemCounts(tx *sqlx.Tx, chainID int64) (*SQLItemCounts, error) {
	counts := SQLItemCounts{}
	err := tx.Get(&counts, `
		SELECT
			COUNT(*) FILTER (WHERE data_type NOT IN ($1, $2)) AS normal_item_count,
			COUNT(*) FILTER (WHERE data_type IN ($1, $2)) AS history_item_count
		FROM entities
		WHERE chain_id = $3 AND deleted = false
	`, HistoryTypeID, HistoryDeleteDirectiveTypeID, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get item counts: %w", err)
	}
	return &counts, nil
}

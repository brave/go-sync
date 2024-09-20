package command

import (
	"context"
	"fmt"

	"github.com/brave/go-sync/cache"
	"github.com/brave/go-sync/datastore"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type ItemCounts struct {
	cache                *cache.Cache
	dynamoDB             datastore.DynamoDatastore
	dynamoItemCounts     *datastore.DynamoItemCounts
	sqlItemCounts        *datastore.SQLItemCounts
	clientID             string
	cacheNewNormalCount  int
	cacheNewHistoryCount int
	sqlTxNewNormalCount  int
	sqlTxNewHistoryCount int
}

func GetItemCounts(cache *cache.Cache, dynamoDB datastore.DynamoDatastore, sqlDB datastore.SQLDatastore, tx *sqlx.Tx, clientID string, chainID int64) (*ItemCounts, error) {
	dynamoItemCounts, err := dynamoDB.GetClientItemCount(clientID)
	if err != nil {
		return nil, err
	}

	sqlItemCounts, err := sqlDB.GetItemCounts(tx, chainID)
	if err != nil {
		return nil, err
	}

	itemCounts := ItemCounts{
		cache:                cache,
		dynamoDB:             dynamoDB,
		dynamoItemCounts:     dynamoItemCounts,
		sqlItemCounts:        sqlItemCounts,
		clientID:             clientID,
		cacheNewNormalCount:  0,
		cacheNewHistoryCount: 0,
		sqlTxNewNormalCount:  0,
		sqlTxNewHistoryCount: 0,
	}
	err = itemCounts.updateInterimItemCounts(false)
	if err != nil {
		return nil, err
	}

	return &itemCounts, nil
}

func (itemCounts *ItemCounts) updateInterimItemCounts(clear bool) error {
	newNormalCount, err := itemCounts.cache.GetInterimCount(context.Background(), itemCounts.clientID, normalCountTypeStr, clear)
	if err != nil {
		return err
	}
	newHistoryCount, err := itemCounts.cache.GetInterimCount(context.Background(), itemCounts.clientID, historyCountTypeStr, clear)
	if err != nil {
		return err
	}
	itemCounts.cacheNewNormalCount = newNormalCount
	itemCounts.cacheNewHistoryCount = newHistoryCount
	return nil
}

func (itemCounts *ItemCounts) RecordChange(dataType int, subtract bool, isStoredInSQL bool) error {
	isHistory := dataType == datastore.HistoryTypeID || dataType == datastore.HistoryDeleteDirectiveTypeID
	if isStoredInSQL {
		delta := 1
		if subtract {
			delta = -1
		}
		if isHistory {
			itemCounts.sqlTxNewHistoryCount += delta
		} else {
			itemCounts.sqlTxNewNormalCount += delta
		}
	} else {
		countType := normalCountTypeStr
		if isHistory {
			countType = historyCountTypeStr
		}
		newCount, err := itemCounts.cache.IncrementInterimCount(context.Background(), itemCounts.clientID, countType, subtract)
		if err != nil {
			return fmt.Errorf("failed to increment history cache count")
		}
		if isHistory {
			itemCounts.cacheNewHistoryCount = newCount
		} else {
			itemCounts.cacheNewNormalCount = newCount
		}
	}
	return nil
}

func (itemCounts *ItemCounts) SumCounts(historyOnly bool) int {
	sum := itemCounts.dynamoItemCounts.SumHistoryCounts() + itemCounts.sqlItemCounts.HistoryItemCount + itemCounts.sqlTxNewHistoryCount + itemCounts.cacheNewHistoryCount
	if !historyOnly {
		sum += itemCounts.dynamoItemCounts.ItemCount + itemCounts.sqlItemCounts.NormalItemCount + itemCounts.sqlTxNewNormalCount + itemCounts.cacheNewNormalCount
	}
	return sum
}

func (itemCounts *ItemCounts) Save() error {
	err := itemCounts.updateInterimItemCounts(true)
	if err != nil {
		return fmt.Errorf("error getting interim item count: %w", err)
	}
	if err = itemCounts.dynamoDB.UpdateClientItemCount(itemCounts.dynamoItemCounts, itemCounts.cacheNewNormalCount, itemCounts.cacheNewHistoryCount); err != nil {
		// We only impose a soft quota limit on the item count for each client, so
		// we only log the error without further actions here. The reason of this
		// is we do not want to pay the cost to ensure strong consistency on this
		// value and we do not want to give up previous DB operations if we cannot
		// update the count this time. In addition, we do not retry this operation
		// either because it is acceptable to miss one time of this update and
		// chances of failing to update the item count multiple times in a row for
		// a single client is quite low.
		log.Error().Err(err).Msg("Update client item count failed")
	}
	return nil
}

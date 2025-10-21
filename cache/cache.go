package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	ttl = 600 * time.Second
)

// Cache is a wrapper for cache accesses.
type Cache struct {
	RedisClient // interface for accessing underlying redis client
}

// NewCache creates a Cache wrapper with underlying redis client assigned.
func NewCache(r RedisClient) *Cache {
	return &Cache{r}
}

// GetTypeMtimeKey returns derived cache key clientID#dataType
func (c *Cache) GetTypeMtimeKey(clientID string, dataType int) string {
	return clientID + "#" + strconv.Itoa(dataType)
}

// GetInterimCountKey returns the key for the entities inserted in
// the DB that were not yet added to the item count
func GetInterimCountKey(clientID string, countType string) string {
	return clientID + "#interim_" + countType
}

// GetAndClearInterimCount returns the amount of entities inserted in
// the DB that were not yet added to the item count
func (c *Cache) GetInterimCount(ctx context.Context, clientID string, countType string, clearCache bool) (int, error) {
	countStr, err := c.Get(ctx, GetInterimCountKey(clientID, countType), clearCache)
	if err != nil {
		log.Error().Err(err).Msg("Get interim count failed")
		return 0, err
	}
	if countStr == "" {
		return 0, nil
	}
	count, err := strconv.Atoi(countStr)
	if err != nil {
		log.Error().Err(err).Msg("Interim count parse failed")
		return 0, err
	}
	return count, nil
}

// IncrementInterimCount increments or decrements the amount of entities inserted in
// the DB that were not yet added to the item count
func (c *Cache) IncrementInterimCount(ctx context.Context, clientID string, countType string, subtract bool) (int, error) {
	return c.Incr(ctx, GetInterimCountKey(clientID, countType), subtract)
}

// SetTypeMtime add an entry into cache where key is clientID#dataType, value
// is the lastest mtime seen on this type for the client.
func (c *Cache) SetTypeMtime(ctx context.Context, clientID string, dataType int, mtime int64) {
	key := c.GetTypeMtimeKey(clientID, dataType)
	val := strconv.FormatInt(mtime, 10)
	err := c.Set(ctx, key, val, ttl)
	if err != nil {
		log.Error().Err(err).Msg("Set value in cache failed")
	}
}

// IsTypeMtimeUpdated check the cache to determine if there might be updates
// for a specific type for a client. It gets the last seen mtime for
// clientID#dataType in the cache, return false if it is found and the value is
// older or equal to client's token, which means the client is already
// up-to-date. In any other cases, it will return false.
func (c *Cache) IsTypeMtimeUpdated(ctx context.Context, clientID string, dataType int, token int64) bool {
	key := c.GetTypeMtimeKey(clientID, dataType)
	cachedTokenStr, err := c.Get(ctx, key, false)

	// If operation failed or cache missed, return true to proceed to querying
	// dynamoDB.
	if err != nil {
		log.Error().Err(err).Msg("Get value from the cache failed")
		return true
	}

	// Cache missed, fall through to querying dynamoDB.
	if cachedTokenStr == "" {
		return true
	}

	// Token parsing is unlikely to fail here, but if it happens, fall through to
	// querying dynamoDB.
	cachedToken, err := strconv.ParseInt(cachedTokenStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Parse cached token value failed")
		return true
	}

	// DB have new updates available since cached mtime is newer.
	if cachedToken > token {
		return true
	}

	// Cached mtime is not newer than client's token, return false to skip
	// querying dynamoDB.
	return false
}

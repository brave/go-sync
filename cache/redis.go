package cache

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient is a wrapper to support clients for standalone redis and redis
// cluster.
type RedisClient interface {
	Set(ctx context.Context, key string, val string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	FlushAll(ctx context.Context) error
}

type redisSimpleClient struct {
	client *redis.Client
}

type redisClusterClient struct {
	client *redis.ClusterClient
}

// NewRedisClient create a client for standalone redis or redis cluster.
func NewRedisClient() RedisClient {
	addrs := strings.Split(os.Getenv("REDIS_URL"), ",")
	cluster := os.Getenv("ENV") != "local"
	poolSize, err := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))
	disableIdentity, _ := strconv.ParseBool(os.Getenv("REDIS_DISABLE_IDENTITY"))

	if err != nil {
		poolSize = 100
	}

	// Fallback to localhost:6397 and non-cluster client if redis env is not set.
	if len(addrs) == 0 {
		addrs = []string{"localhost:6397"}
		cluster = false
	}

	var r RedisClient

	if !cluster {
		client := redis.NewClient(&redis.Options{
			Addr:             addrs[0],
			DisableIndentity: disableIdentity,
		})
		r = &redisSimpleClient{client}
	} else {
		client := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:            addrs,
			PoolSize:         poolSize,
			ReadOnly:         true,
			DisableIndentity: disableIdentity,
		})
		r = &redisClusterClient{client}
	}

	return r
}

func (r *redisSimpleClient) Set(ctx context.Context, key string, val string, ttl time.Duration) error {
	return r.client.Set(ctx, key, val, ttl).Err()
}

func (r *redisSimpleClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *redisSimpleClient) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *redisSimpleClient) FlushAll(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

func (r *redisClusterClient) Set(ctx context.Context, key string, val string, ttl time.Duration) error {
	return r.client.Set(ctx, key, val, ttl).Err()
}

func (r *redisClusterClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *redisClusterClient) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *redisClusterClient) FlushAll(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

package cache

import (
	"context"
	"fmt"
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
	Incr(ctx context.Context, key string, subtract bool) (int, error)
	Get(ctx context.Context, key string, delete bool) (string, error)
	Del(ctx context.Context, keys ...string) error
	FlushAll(ctx context.Context) error
	SubscribeAndWait(ctx context.Context, channel string) error
}

type redisClientImpl struct {
	client redis.UniversalClient
}

// NewRedisClient create a client for standalone redis or redis cluster.
func NewRedisClient() RedisClient {
	addrs := strings.Split(os.Getenv("REDIS_URL"), ",")
	env := os.Getenv("ENV")
	cluster := env != "local" && env != ""
	poolSize, err := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))
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
			Addr: addrs[0],
		})
		r = &redisClientImpl{client}
	} else {
		client := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,
			PoolSize: poolSize,
			ReadOnly: true,
		})
		r = &redisClientImpl{client}
	}

	return r
}

func (r *redisClientImpl) Set(ctx context.Context, key string, val string, ttl time.Duration) error {
	return r.client.Set(ctx, key, val, ttl).Err()
}

func (r *redisClientImpl) Incr(ctx context.Context, key string, subtract bool) (int, error) {
	var res *redis.IntCmd
	if subtract {
		res = r.client.Decr(ctx, key)
	} else {
		res = r.client.Incr(ctx, key)
	}
	val, err := res.Result()
	return int(val), err
}

func (r *redisClientImpl) Get(ctx context.Context, key string, delete bool) (string, error) {
	var res *redis.StringCmd
	if delete {
		res = r.client.GetDel(ctx, key)
	} else {
		res = r.client.Get(ctx, key)
	}
	val, err := res.Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *redisClientImpl) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *redisClientImpl) FlushAll(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

func (r *redisClientImpl) SubscribeAndWait(ctx context.Context, channel string) error {
	pubsub := r.client.Subscribe(ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return fmt.Errorf("redis channel unexpectedly closed")
			}
			if msg != nil {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

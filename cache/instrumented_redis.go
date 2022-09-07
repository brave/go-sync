package cache

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../.prom-gowrap.tmpl template

//go:generate gowrap gen -p github.com/odedlaz/go-sync/cache -i RedisClient -t ../.prom-gowrap.tmpl -o instrumented_redis.go

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// RedisClientWithPrometheus implements RedisClient interface with all methods wrapped
// with Prometheus metrics
type RedisClientWithPrometheus struct {
	base         RedisClient
	instanceName string
}

var redisclientDurationSummaryVec = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "redisclient_duration_seconds",
		Help:       "redisclient runtime duration and result",
		MaxAge:     time.Minute,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"instance_name", "method", "result"})

// NewRedisClientWithPrometheus returns an instance of the RedisClient decorated with prometheus summary metric
func NewRedisClientWithPrometheus(base RedisClient, instanceName string) RedisClientWithPrometheus {
	return RedisClientWithPrometheus{
		base:         base,
		instanceName: instanceName,
	}
}

// FlushAll implements RedisClient
func (_d RedisClientWithPrometheus) FlushAll(ctx context.Context) (err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		redisclientDurationSummaryVec.WithLabelValues(_d.instanceName, "FlushAll", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.FlushAll(ctx)
}

// Get implements RedisClient
func (_d RedisClientWithPrometheus) Get(ctx context.Context, key string) (s1 string, err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		redisclientDurationSummaryVec.WithLabelValues(_d.instanceName, "Get", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.Get(ctx, key)
}

// Set implements RedisClient
func (_d RedisClientWithPrometheus) Set(ctx context.Context, key string, val string, ttl time.Duration) (err error) {
	_since := time.Now()
	defer func() {
		result := "ok"
		if err != nil {
			result = "error"
		}

		redisclientDurationSummaryVec.WithLabelValues(_d.instanceName, "Set", result).Observe(time.Since(_since).Seconds())
	}()
	return _d.base.Set(ctx, key, val, ttl)
}

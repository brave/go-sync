package datastore

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ratelimit"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	// Strings for the primary key
	pk     string = "ClientID"
	sk     string = "ID"
	projPk string = "ClientID, ID"

	// Strings for (ClientID, DataTypeMtime) GSI
	clientIDDataTypeMtimeIdx   string = "ClientIDDataTypeMtimeIndex"
	clientIDDataTypeMtimeIdxPk string = "ClientID"
	clientIDDataTypeMtimeIdxSk string = "DataTypeMtime"

	// Default retry configuration for DynamoDB API calls via the AWS SDK.
	// Each value can be overridden at runtime via the corresponding env var.
	//
	// defaultRetryMaxAttempts: total attempts per operation (1 initial + 4 retries).
	// The SDK default is 3. We use 5 because DynamoDB throttles are common under
	// bursty sync traffic and the extra attempts, combined with jittered backoff,
	// let transient 5xx / throttle errors resolve without surfacing to callers.
	// Override: DYNAMO_RETRY_MAX_ATTEMPTS.
	defaultRetryMaxAttempts = 5

	// defaultRetryMaxBackoff: upper bound on the jittered exponential backoff
	// between retries. The SDK default is 20s, which is too long for an
	// interactive sync path. 5s keeps worst-case total retry time bounded while
	// still giving DynamoDB enough breathing room to shed load.
	// Override: DYNAMO_RETRY_MAX_BACKOFF (Go duration string, e.g. "5s").
	defaultRetryMaxBackoff = 5 * time.Second

	// defaultRetryTokenBucketSize: capacity of the client-side token-bucket rate
	// limiter. Each retry costs tokens; when the bucket is exhausted the SDK
	// fails the operation with a QuotaExceededError.
	// Override: DYNAMO_RETRY_TOKEN_BUCKET_SIZE.
	defaultRetryTokenBucketSize = 1000
)

var (
	// Table is the name of the table in dynamoDB, could be modified in tests.
	Table = os.Getenv("TABLE_NAME")
)

// PrimaryKey struct is used to represent the primary key of our table.
type PrimaryKey struct {
	ClientID string // Hash key
	ID       string // Range key
}

// Dynamo is a Datastore wrapper around a dynamoDB.
type Dynamo struct {
	*dynamodb.Client
}

func envInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultVal
}

func envDuration(key string, defaultVal time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return defaultVal
}

// newDynamoHTTPClient returns an HTTP client with phase-specific timeouts
// tuned for in-region DynamoDB calls.
//
// Unlike http.Client.Timeout (a blanket timeout covering the entire request
// lifecycle), this configures each phase independently:
//   - Connect: 1.1s — matching SDK DefaultsModeInRegion, fails fast on
//     unreachable endpoints instead of hanging for 30s.
//   - TLS handshake: 1.1s — same rationale; TLS to DynamoDB in-region is <50ms
//     at p99, so 1.1s is generous.
//   - ResponseHeaderTimeout: 5s — bounds time waiting for DynamoDB to start
//     sending a response after the request is fully sent. DynamoDB p99 is
//     single-digit ms; 5s catches stalled connections without killing slow but
//     progressing requests. The SDK has no default for this, leaving a gap
//     where a stalled response hangs until context deadline (60s).
//   - http.Client.Timeout: 30s — last-resort safety net. Some write paths use
//     context.Background() with no deadline; without this, a fully stalled
//     connection that passes all phase timeouts could hang indefinitely. At
//     the http.Client level a timeout produces a retryable error (Timeout()
//     == true), so the SDK retries the individual HTTP request rather than
//     killing the entire operation.
func newDynamoHTTPClient() *http.Client {
	const connectTimeout = 1100 * time.Millisecond

	dialer := &net.Dialer{
		Timeout:   connectTimeout,
		KeepAlive: 30 * time.Second,
	}
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        envInt("DYNAMO_MAX_IDLE_CONNS", 600),
			MaxIdleConnsPerHost: envInt("DYNAMO_MAX_IDLE_CONNS_PER_HOST", 400),

			DialContext:           dialer.DialContext,
			TLSHandshakeTimeout:   connectTimeout,
			ResponseHeaderTimeout: 5 * time.Second,
			IdleConnTimeout:       90 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}
}

// NewDynamo returns a dynamoDB client to be used.
func NewDynamo() (*Dynamo, error) {
	ctx := context.Background()

	// Load default AWS configuration
	configOpts := []func(*config.LoadOptions) error{
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithHTTPClient(newDynamoHTTPClient()),
	}
	if os.Getenv("LOG_SDK_RETRIES") == "1" {
		configOpts = append(configOpts, config.WithClientLogMode(aws.LogRetries))
	}
	cfg, err := config.LoadDefaultConfig(ctx, configOpts...)
	if err != nil {
		return nil, fmt.Errorf("error loading AWS config: %w", err)
	}

	// Create DynamoDB client with optional endpoint override
	endpoint := os.Getenv("AWS_ENDPOINT")
	db := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}

		o.Retryer = retry.NewStandard(func(so *retry.StandardOptions) {
			so.MaxAttempts = envInt("DYNAMO_RETRY_MAX_ATTEMPTS", defaultRetryMaxAttempts)
			so.MaxBackoff = envDuration("DYNAMO_RETRY_MAX_BACKOFF", defaultRetryMaxBackoff)
			so.RateLimiter = ratelimit.NewTokenRateLimit(
				uint(envInt("DYNAMO_RETRY_TOKEN_BUCKET_SIZE", defaultRetryTokenBucketSize)),
			)
		})
	})

	return &Dynamo{db}, nil
}

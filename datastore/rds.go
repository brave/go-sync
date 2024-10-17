package datastore

import (
	"context"
	"database/sql/driver"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/jackc/pgx/stdlib"
)

const defaultRegion = "us-west-2"

const (
	rdsPortKey   = "RDS_DATABASE_PORT"
	rdsHostKey   = "RDS_WRITER_ENDPOINT"
	rdsUserKey   = "RDS_USER"
	rdsDbNameKey = "RDS_DATABASE_NAME"
	regionKey    = "AWS_REGION"
)

type rdsConnector struct {
	hostAndPort    string
	dbName         string
	user           string
	token          string
	region         string
	tokenCacheTime time.Time
	mu             sync.Mutex
}

func newRDSConnector() *rdsConnector {
	port := os.Getenv(rdsPortKey)
	host := os.Getenv(rdsHostKey)
	user := os.Getenv(rdsUserKey)
	dbName := os.Getenv(rdsDbNameKey)
	region := os.Getenv(regionKey)

	if region == "" {
		region = defaultRegion
	}
	hostAndPort := fmt.Sprintf("%s:%s", host, port)
	return &rdsConnector{
		hostAndPort: hostAndPort,
		dbName:      dbName,
		user:        user,
		region:      region,
	}
}

func (c *rdsConnector) getConnectionString(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Since(c.tokenCacheTime) > 10*time.Minute {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to load AWS config")
		}

		token, err := auth.BuildAuthToken(
			ctx, c.hostAndPort, c.region, c.user, cfg.Credentials)
		if err != nil {
			return "", fmt.Errorf("failed to create authentication token: %w", err)
		}
		c.token = url.QueryEscape(token)
		c.tokenCacheTime = time.Now()
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=require", c.user, c.token, c.hostAndPort, c.dbName), nil
}

func (c *rdsConnector) Connect(ctx context.Context) (driver.Conn, error) {
	connStr, err := c.getConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	return stdlib.GetDefaultDriver().Open(connStr)
}

func (c *rdsConnector) Driver() driver.Driver {
	return c
}

func (c *rdsConnector) Open(_ string) (driver.Conn, error) {
	return nil, fmt.Errorf("open method unsupported")
}

package datastore

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
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

// NewDynamo returns a dynamoDB client to be used.
func NewDynamo() (*Dynamo, error) {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 50,
		},
	}

	ctx := context.Background()

	// Load default AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("error loading AWS config: %w", err)
	}

	// Create DynamoDB client with optional endpoint override
	endpoint := os.Getenv("AWS_ENDPOINT")
	db := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})

	return &Dynamo{db}, nil
}

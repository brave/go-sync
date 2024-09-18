package datastore

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
	Table               = os.Getenv("TABLE_NAME")
	defaultTestEndpoint = "http://localhost:8000"
	defaultTestRegion   = "us-west-2"
)

// PrimaryKey struct is used to represent the primary key of our table.
type PrimaryKey struct {
	ClientID string // Hash key
	ID       string // Range key
}

// Dynamo is a Datastore wrapper around a dynamoDB.
type Dynamo struct {
	*dynamodb.DynamoDB
}

// NewDynamo returns a dynamoDB client to be used.
func NewDynamo(isTesting bool) (*Dynamo, error) {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 50,
		},
	}

	endpoint := os.Getenv("AWS_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	if endpoint == "" && region == "" && isTesting {
		endpoint = defaultTestEndpoint
		region = defaultTestRegion
	}

	awsConfig := aws.NewConfig().WithRegion(region).WithEndpoint(endpoint).WithHTTPClient(httpClient)

	if isTesting {
		awsConfig = awsConfig.WithCredentials(credentials.NewStaticCredentials("GOSYNC", "GOSYNC", "GOSYNC"))
	}

	sess, err := session.NewSession(awsConfig)

	if err != nil {
		return nil, fmt.Errorf("error creating new AWS session: %w", err)
	}

	db := dynamodb.New(sess)
	return &Dynamo{db}, nil
}

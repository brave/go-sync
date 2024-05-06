package datastore

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
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
	Table = os.Getenv("TABLE_NAME")
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
func NewDynamo() (*Dynamo, error) {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 50,
		},
	}

	awsConfig := aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")).WithEndpoint(os.Getenv("AWS_ENDPOINT")).WithHTTPClient(httpClient)
	sess, err := session.NewSession(awsConfig)

	if err != nil {
		return nil, fmt.Errorf("error creating new AWS session: %w", err)
	}

	db := dynamodb.New(sess)
	return &Dynamo{db}, nil
}

package datastoretest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/brave/go-sync/datastore"
)

// DeleteTable deletes datastore.Table in dynamoDB.
func DeleteTable(dynamo *datastore.Dynamo) error {
	_, err := dynamo.DeleteTable(context.Background(),
		&dynamodb.DeleteTableInput{TableName: aws.String(datastore.Table)})
	if err != nil {
		var notFoundException *types.ResourceNotFoundException
		if errors.As(err, &notFoundException) {
			// Return as successful if the table is not existed.
			return nil
		}
		return fmt.Errorf("error deleting table: %w", err)
	}

	// Wait for table to be deleted using waiter
	waiter := dynamodb.NewTableNotExistsWaiter(dynamo)
	return waiter.Wait(context.Background(),
		&dynamodb.DescribeTableInput{TableName: aws.String(datastore.Table)},
		5*time.Minute)
}

// CreateTable creates datastore.Table in dynamoDB.
func CreateTable(dynamo *datastore.Dynamo) error {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../../")
	raw, err := os.ReadFile(filepath.Join(root, "schema/dynamodb/table.json"))
	if err != nil {
		return fmt.Errorf("error reading table.json: %w", err)
	}

	var input dynamodb.CreateTableInput
	err = json.Unmarshal(raw, &input)
	if err != nil {
		return fmt.Errorf("error unmarshalling raw data from table.json: %w", err)
	}
	input.TableName = aws.String(datastore.Table)

	_, err = dynamo.CreateTable(context.Background(), &input)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	// Wait for table to be active using waiter
	waiter := dynamodb.NewTableExistsWaiter(dynamo)
	return waiter.Wait(context.Background(),
		&dynamodb.DescribeTableInput{TableName: aws.String(datastore.Table)},
		5*time.Minute)
}

// ResetTable deletes and creates datastore.Table in dynamoDB.
func ResetTable(dynamo *datastore.Dynamo) error {
	if err := DeleteTable(dynamo); err != nil {
		return fmt.Errorf("error deleting table to reset table: %w", err)
	}
	return CreateTable(dynamo)
}

// ScanSyncEntities scans the dynamoDB table and returns all sync items.
func ScanSyncEntities(dynamo *datastore.Dynamo) ([]datastore.SyncEntity, error) {
	filter := expression.AttributeExists(expression.Name("Version"))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return nil, fmt.Errorf("error building expression to scan sync entitites: %w", err)
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(datastore.Table),
	}
	out, err := dynamo.Scan(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("error doing scan for sync entities: %w", err)
	}
	syncItems := []datastore.SyncEntity{}
	err = attributevalue.UnmarshalListOfMaps(out.Items, &syncItems)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling sync entitites: %w", err)
	}

	return syncItems, nil
}

// ScanTagItems scans the dynamoDB table and returns all tag items.
func ScanTagItems(dynamo *datastore.Dynamo) ([]datastore.ServerClientUniqueTagItem, error) {
	filter := expression.And(
		expression.AttributeNotExists(expression.Name("ExpireAt")),
		expression.AttributeNotExists(expression.Name("Version")))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return nil, fmt.Errorf("error building expression to scan tag items: %w", err)
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(datastore.Table),
	}
	out, err := dynamo.Scan(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("error doing scan for tag items: %w", err)
	}
	tagItems := []datastore.ServerClientUniqueTagItem{}
	err = attributevalue.UnmarshalListOfMaps(out.Items, &tagItems)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling tag items: %w", err)
	}

	return tagItems, nil
}

// ScanClientItemCounts scans the dynamoDB table and returns all client item
// counts.
func ScanClientItemCounts(dynamo *datastore.Dynamo) ([]datastore.ClientItemCounts, error) {
	filter := expression.AttributeExists(expression.Name("ItemCount"))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return nil, fmt.Errorf("error building expression to scan item counts: %w", err)
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(datastore.Table),
	}
	out, err := dynamo.Scan(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("error doing scan for item counts: %w", err)
	}
	clientItemCounts := []datastore.ClientItemCounts{}
	err = attributevalue.UnmarshalListOfMaps(out.Items, &clientItemCounts)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling item counts: %w", err)
	}

	return clientItemCounts, nil
}

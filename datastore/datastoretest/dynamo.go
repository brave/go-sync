package datastoretest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/brave/go-sync/datastore"
)

// DeleteTable deletes datastore.Table in dynamoDB.
func DeleteTable(dynamo *datastore.Dynamo) error {
	_, err := dynamo.DeleteTable(
		&dynamodb.DeleteTableInput{TableName: aws.String(datastore.Table)})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			// Return as successful if the table is not existed.
			if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
				return nil
			}
			return err
		} else {
			return fmt.Errorf("error deleting table: %w", err)
		}
	}

	return dynamo.WaitUntilTableNotExists(
		&dynamodb.DescribeTableInput{TableName: aws.String(datastore.Table)})
}

// CreateTable creates datastore.Table in dynamoDB.
func CreateTable(dynamo *datastore.Dynamo) error {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../../")
	raw, err := ioutil.ReadFile(filepath.Join(root, "schema/dynamodb/table.json"))
	if err != nil {
		return fmt.Errorf("error reading table.json: %w", err)
	}

	var input dynamodb.CreateTableInput
	err = json.Unmarshal(raw, &input)
	if err != nil {
		return fmt.Errorf("error unmarshalling raw data from table.json: %w", err)
	}
	input.TableName = aws.String(datastore.Table)

	_, err = dynamo.CreateTable(&input)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	return dynamo.WaitUntilTableExists(
		&dynamodb.DescribeTableInput{TableName: aws.String(datastore.Table)})
}

// ResetDynamoTable deletes and creates datastore.Table in dynamoDB.
func ResetDynamoTable(dynamo *datastore.Dynamo) error {
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
	out, err := dynamo.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("error doing scan for sync entities: %w", err)
	}
	syncItems := []datastore.SyncEntity{}
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &syncItems)
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
	out, err := dynamo.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("error doing scan for tag items: %w", err)
	}
	tagItems := []datastore.ServerClientUniqueTagItem{}
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &tagItems)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling tag items: %w", err)
	}

	return tagItems, nil
}

// ScanClientItemCounts scans the dynamoDB table and returns all client item
// counts.
func ScanClientItemCounts(dynamo *datastore.Dynamo) ([]datastore.DynamoItemCounts, error) {
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
	out, err := dynamo.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("error doing scan for item counts: %w", err)
	}
	clientItemCounts := []datastore.DynamoItemCounts{}
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &clientItemCounts)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling item counts: %w", err)
	}

	return clientItemCounts, nil
}

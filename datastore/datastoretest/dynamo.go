package datastoretest

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/brave-experiments/sync-server/datastore"
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
		} else {
			return err
		}
	}

	return dynamo.WaitUntilTableNotExists(
		&dynamodb.DescribeTableInput{TableName: aws.String(datastore.Table)})
}

// CreateTable creates datastore.Table in dynamoDB.
func CreateTable(dynamo *datastore.Dynamo) error {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../../")
	raw, err := ioutil.ReadFile(filepath.Join(root, "dynamo_local/table.json"))
	if err != nil {
		return err
	}

	var input dynamodb.CreateTableInput
	err = json.Unmarshal(raw, &input)
	if err != nil {
		return err
	}
	input.TableName = aws.String(datastore.Table)

	_, err = dynamo.CreateTable(&input)
	if err != nil {
		return err
	}

	return dynamo.WaitUntilTableExists(
		&dynamodb.DescribeTableInput{TableName: aws.String(datastore.Table)})
}

// ResetTable deletes and creates datastore.Table in dynamoDB.
func ResetTable(dynamo *datastore.Dynamo) error {
	if err := DeleteTable(dynamo); err != nil {
		return err
	}
	return CreateTable(dynamo)
}

// ScanClientTokens scans the dynamoDB table and returns all client-token items.
func ScanClientTokens(dynamo *datastore.Dynamo) ([]datastore.ClientToken, error) {
	filter := expression.AttributeExists(expression.Name("ExpireAt"))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(datastore.Table),
	}
	out, err := dynamo.Scan(input)
	if err != nil {
		return nil, err
	}
	clientTokens := []datastore.ClientToken{}
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &clientTokens)
	if err != nil {
		return nil, err
	}

	return clientTokens, nil
}

// ScanSyncEntities scans the dynamoDB table and returns all sync items.
func ScanSyncEntities(dynamo *datastore.Dynamo) ([]datastore.SyncEntity, error) {
	filter := expression.AttributeExists(expression.Name("Version"))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(datastore.Table),
	}
	out, err := dynamo.Scan(input)
	if err != nil {
		return nil, err
	}
	syncItems := []datastore.SyncEntity{}
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &syncItems)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(datastore.Table),
	}
	out, err := dynamo.Scan(input)
	if err != nil {
		return nil, err
	}
	tagItems := []datastore.ServerClientUniqueTagItem{}
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &tagItems)
	if err != nil {
		return nil, err
	}

	return tagItems, nil
}

// ScanClientItemCounts scans the dynamoDB table and returns all client item
// counts.
func ScanClientItemCounts(dynamo *datastore.Dynamo) ([]datastore.ClientItemCount, error) {
	filter := expression.AttributeExists(expression.Name("ItemCount"))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(datastore.Table),
	}
	out, err := dynamo.Scan(input)
	if err != nil {
		return nil, err
	}
	clientItemCounts := []datastore.ClientItemCount{}
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &clientItemCounts)
	if err != nil {
		return nil, err
	}

	return clientItemCounts, nil
}

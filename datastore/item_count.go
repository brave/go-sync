package datastore

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	itemCountAttrName string = "ItemCount"
)

// ClientItemCount is used to marshal and unmarshal ClientItemCount items in
// dynamoDB.
type ClientItemCount struct {
	ClientID  string
	ID        string
	ItemCount int
}

// ClientItemCountByClientID  implements sort.Interface for []ClientItemCount
// based on ClientID.
type ClientItemCountByClientID []ClientItemCount

func (a ClientItemCountByClientID) Len() int      { return len(a) }
func (a ClientItemCountByClientID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ClientItemCountByClientID) Less(i, j int) bool {
	return a[i].ClientID < a[j].ClientID
}

// GetClientItemCount returns the count of non-deleted sync items stored for
// a given client.
func (dynamo *Dynamo) GetClientItemCount(clientID string) (int, error) {
	primaryKey := PrimaryKey{ClientID: clientID, ID: clientID}
	key, err := dynamodbattribute.MarshalMap(primaryKey)
	if err != nil {
		return 0, fmt.Errorf("error marshalling primary key to get item-count item: %w", err)
	}

	input := &dynamodb.GetItemInput{
		Key:                  key,
		ProjectionExpression: aws.String(itemCountAttrName),
		TableName:            aws.String(Table),
	}

	out, err := dynamo.GetItem(input)
	if err != nil {
		return 0, fmt.Errorf("error getting an item-count item: %w", err)
	}

	clientItemCount := &ClientItemCount{}
	err = dynamodbattribute.UnmarshalMap(out.Item, clientItemCount)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling item-count item: %w", err)
	}
	return clientItemCount.ItemCount, nil
}

// UpdateClientItemCount updates the count of non-deleted sync items for a
// given client stored in the dynamoDB.
func (dynamo *Dynamo) UpdateClientItemCount(clientID string, count int) error {
	primaryKey := PrimaryKey{ClientID: clientID, ID: clientID}
	key, err := dynamodbattribute.MarshalMap(primaryKey)
	if err != nil {
		return fmt.Errorf("error marshalling primary key to update item-count item: %w", err)
	}

	update := expression.Add(expression.Name(itemCountAttrName), expression.Value(count))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("error building expression to update item-count item: %w", err)
	}

	input := &dynamodb.UpdateItemInput{
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		TableName:                 aws.String(Table),
	}

	_, err = dynamo.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("error updating item-count item in dynamoDB: %w", err)
	}
	return nil
}

package datastore

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/brave/go-sync/utils"
)

// ClientToken is a struct used to marshal and unmarshal client_token items in
// dynamoDB.
type ClientToken struct {
	ClientID string
	Token    string `dynamodbav:"ID"`
	ExpireAt int64
}

// InsertClientToken insert a new client_token item into dynamoDB.
func (dynamo *Dynamo) InsertClientToken(id string, token string, expireAt int64) error {
	clientToken := ClientToken{ClientID: id, Token: token, ExpireAt: expireAt}
	av, err := dynamodbattribute.MarshalMap(clientToken)
	if err != nil {
		return fmt.Errorf("error marshalling client token item to insert client token: %w", err)
	}
	cond := expression.AttributeNotExists(expression.Name(pk))
	expr, err := expression.NewBuilder().WithCondition(cond).Build()
	if err != nil {
		return fmt.Errorf("error building expression to insert client token: %w", err)
	}
	input := &dynamodb.PutItemInput{
		Item:                      av,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
		TableName:                 aws.String(Table),
	}
	_, err = dynamo.PutItem(input)
	if err != nil {
		return fmt.Errorf("error calling PutItem to insert client token: %w", err)
	}
	return nil
}

// GetClientID queries the table using (ID, ExpireAt) GSI and returns the
// clientID if there is a corresponding non-expired token existed in the table.
func (dynamo *Dynamo) GetClientID(token string) (string, error) {
	pkCond := expression.Key(idExpireAtIdxPk).Equal(expression.Value(token))
	skCond := expression.Key(idExpireAtIdxSk).GreaterThan(expression.Value(utils.UnixMilli(time.Now())))
	keyCond := expression.KeyAnd(pkCond, skCond)
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return "", fmt.Errorf("error building expression to get client ID: %w", err)
	}
	input := &dynamodb.QueryInput{
		IndexName:                 aws.String(idExpireAtIdx),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String(Table),
	}
	out, err := dynamo.Query(input)
	if err != nil {
		return "", fmt.Errorf("error doing query to get client ID: %w", err)
	}
	if *(out.Count) == 0 {
		return "", fmt.Errorf("Not a valid token")
	}

	clientTokens := []ClientToken{}
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &clientTokens)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling client tokens: %w", err)
	}

	return clientTokens[0].ClientID, nil
}

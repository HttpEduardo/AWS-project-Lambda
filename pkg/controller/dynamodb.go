package controller

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (ctrl *Controller) existUserNameAtDynamoDB(ctx context.Context, svc *dynamodb.DynamoDB, userName string) (bool, error) {
	result, err := svc.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"UserName": {
				S: aws.String(userName),
			},
		},
		TableName: aws.String(ctrl.Config.DynamoDBTableName),
	})
	if err != nil {
		return false, fmt.Errorf("get an item from DynamoDB: %w", err)
	}
	return len(result.Item) != 0, nil
}

func (ctrl *Controller) addUserNameToDynamoDB(ctx context.Context, svc *dynamodb.DynamoDB, userName, expiredAt string) error {
	if _, err := svc.UpdateItemWithContext(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(ctrl.Config.DynamoDBTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserName": {
				S: aws.String(userName),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":ExpiredAt": {
				N: aws.String(expiredAt),
			},
		},
		UpdateExpression: aws.String("SET ExpiredAt = :ExpiredAt"),
	}); err != nil {
		return fmt.Errorf("update an item of DynamoDB: %w", err)
	}
	return nil
}

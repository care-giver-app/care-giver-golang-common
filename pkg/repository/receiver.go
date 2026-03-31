package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/care-giver-app/care-giver-golang-common/pkg/log"
	"github.com/care-giver-app/care-giver-golang-common/pkg/receiver"
	"go.uber.org/zap"
)

const (
	receiverID = "receiver_id"
)

type ReceiverRepositoryProvider interface {
	CreateReceiver(r receiver.Receiver) error
	GetReceiver(rid string) (receiver.Receiver, error)
}

type ReceiverRepository struct {
	Ctx       context.Context
	Client    DynamodbClientProvider
	TableName string
	logger    *zap.Logger
}

func NewReceiverRespository(ctx context.Context, tableName string, client DynamodbClientProvider, logger *zap.Logger) *ReceiverRepository {
	return &ReceiverRepository{
		Ctx:       ctx,
		Client:    client,
		TableName: tableName,
		logger:    logger.With(zap.String(log.TableNameLogKey, tableName)),
	}
}

func (rr *ReceiverRepository) CreateReceiver(r receiver.Receiver) error {
	rr.logger.Info("adding receiver to db", zap.Any(log.ReceiverIDLogKey, r.ReceiverID))

	rr.logger.Info("marshalling receiver struct")
	av, err := attributevalue.MarshalMap(r)
	if err != nil {
		return err
	}

	rr.logger.Info("inserting item into db", zap.Any("item", av))
	_, err = rr.Client.PutItem(rr.Ctx, &dynamodb.PutItemInput{
		TableName: aws.String(rr.TableName),
		Item:      av,
	})
	if err != nil {
		return err
	}
	rr.logger.Info("successfully inserted item")

	return nil
}

func (rr *ReceiverRepository) GetReceiver(rid string) (receiver.Receiver, error) {
	rr.logger.Info("getting receiver from db", zap.Any(log.ReceiverIDLogKey, rid))
	result, err := rr.Client.GetItem(rr.Ctx, &dynamodb.GetItemInput{
		TableName: &rr.TableName,
		Key: map[string]types.AttributeValue{
			receiverID: &types.AttributeValueMemberS{Value: rid},
		},
	})
	if err != nil {
		return receiver.Receiver{}, err
	}

	var r receiver.Receiver
	err = attributevalue.UnmarshalMap(result.Item, &r)
	if err != nil {
		return receiver.Receiver{}, err
	}

	return r, nil
}

func (rr *ReceiverRepository) UpdateReceiver(r receiver.Receiver) error {
	rr.logger.Info("updating receiver in db", zap.String(log.ReceiverIDLogKey, r.ReceiverID))

	if r.ReceiverID == "" {
		return fmt.Errorf("receiver id is required")
	}

	// Marshal the struct to get attribute values, excluding the key
	av, err := attributevalue.MarshalMap(r)
	if err != nil {
		return err
	}

	// Remove the primary key from the update attributes
	delete(av, receiverID)

	// If no fields to update, return early
	if len(av) == 0 {
		rr.logger.Info("no fields to update")
		return nil
	}

	// Build the SET expressions dynamically
	var setExpressions []string
	expressionAttributeNames := make(map[string]string)
	expressionAttributeValues := make(map[string]types.AttributeValue)

	for i, attrName := range av {
		attrValue := av[i]
		placeholderName := fmt.Sprintf("#attr%d", i)
		placeholderValue := fmt.Sprintf(":val%d", i)
		setExpressions = append(setExpressions, fmt.Sprintf("%s = %s", placeholderName, placeholderValue))
		expressionAttributeNames[placeholderName] = attrName
		expressionAttributeValues[placeholderValue] = attrValue
	}

	updateExpression := "SET " + strings.Join(setExpressions, ", ")

	_, err = rr.Client.UpdateItem(rr.Ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(rr.TableName),
		Key: map[string]types.AttributeValue{
			receiverID: &types.AttributeValueMemberS{Value: r.ReceiverID},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	})

	if err != nil {
		return err
	}

	rr.logger.Info("successfully updated receiver")
	return nil
}

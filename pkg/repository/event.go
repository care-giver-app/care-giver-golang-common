package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/care-giver-app/care-giver-golang-common/pkg/event"
	"github.com/care-giver-app/care-giver-golang-common/pkg/log"
	"go.uber.org/zap"
)

type EventRepositoryProvider interface {
	AddEvent(e *event.Entry) error
	GetEvents(rid string) ([]event.Entry, error)
	DeleteEvent(rid, eid string) error
}

type EventRepository struct {
	Ctx       context.Context
	Client    DynamodbClientProvider
	TableName string
	logger    *zap.Logger
}

func NewEventRespository(ctx context.Context, tableName string, client DynamodbClientProvider, logger *zap.Logger) *EventRepository {
	return &EventRepository{
		Ctx:       ctx,
		Client:    client,
		TableName: tableName,
		logger:    logger.With(zap.String(log.TableNameLogKey, tableName)),
	}
}

func (er *EventRepository) AddEvent(e *event.Entry) error {
	er.logger.Info("adding receiver event to db")

	er.logger.Info("marshalling receiver event struct")
	av, err := attributevalue.MarshalMap(e)
	if err != nil {
		return err
	}

	er.logger.Info("inserting item into db", zap.Any("item", av))
	_, err = er.Client.PutItem(er.Ctx, &dynamodb.PutItemInput{
		TableName: aws.String(er.TableName),
		Item:      av,
	})
	if err != nil {
		return err
	}
	er.logger.Info("successfully inserted item")

	return nil
}

type TimestampBound struct {
	Lower string
	Upper string
}

func (er *EventRepository) GetEvents(rid string, bound TimestampBound) ([]event.Entry, error) {
	er.logger.Info("retrieving receiver events from db", zap.String(log.ReceiverIDLogKey, string(rid)))

	keyCondition := "receiver_id = :rid"
	expressionAttributeValues := map[string]types.AttributeValue{
		":rid": &types.AttributeValueMemberS{Value: string(rid)},
	}

	if bound.Upper != "" && bound.Lower != "" {
		keyCondition = fmt.Sprintf("%s %s", keyCondition, "AND timestamp BETWEEN :timelower AND :timeupper")
		expressionAttributeValues[":timelower"] = &types.AttributeValueMemberS{Value: bound.Lower}
		expressionAttributeValues[":timeupper"] = &types.AttributeValueMemberS{Value: bound.Upper}
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(er.TableName),
		IndexName:                 aws.String("receiver-timestamp"),
		KeyConditionExpression:    aws.String(keyCondition),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	result, err := er.Client.Query(er.Ctx, queryInput)
	if err != nil {
		return nil, err
	}

	var eventsList []event.Entry
	err = attributevalue.UnmarshalListOfMaps(result.Items, &eventsList)
	if err != nil {
		er.logger.Error("error unmarshalling events list", zap.Error(err))
		return nil, err
	}

	return eventsList, nil
}

func (er *EventRepository) DeleteEvent(rid, eid string) error {
	er.logger.Info("deleting receiver event from db", zap.String(log.EventIDLogKey, eid))

	_, err := er.Client.DeleteItem(er.Ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(er.TableName),
		Key: map[string]types.AttributeValue{
			"receiver_id": &types.AttributeValueMemberS{Value: rid},
			"event_id":    &types.AttributeValueMemberS{Value: eid},
		},
	})

	if err != nil {
		return err
	}

	er.logger.Info("successfully deleted event")
	return nil
}

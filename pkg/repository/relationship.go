package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/care-giver-app/care-giver-golang-common/pkg/log"
	"github.com/care-giver-app/care-giver-golang-common/pkg/relationship"
	"go.uber.org/zap"
)

type RelationshipRepositoryProvider interface {
	AddRelationship(r *relationship.Relationship) error
	GetRelationship(userID string, receiverID string) (*relationship.Relationship, error)
	GetRelationshipsByUser(userID string) ([]relationship.Relationship, error)
	DeleteRelationship(userID string, receiverID string) error
	GetRelationshipsByEmailNotifications() ([]relationship.Relationship, error)
}

type RelationshipRepository struct {
	Ctx       context.Context
	Client    DynamodbClientProvider
	TableName string
	logger    *zap.Logger
}

func NewRelationshipRepository(ctx context.Context, tableName string, client DynamodbClientProvider, logger *zap.Logger) *RelationshipRepository {
	return &RelationshipRepository{
		Ctx:       ctx,
		Client:    client,
		TableName: tableName,
		logger:    logger.With(zap.String(log.TableNameLogKey, tableName)),
	}
}

func (rr *RelationshipRepository) AddRelationship(r *relationship.Relationship) error {
	rr.logger.Info("adding user receiver relationship to db")

	rr.logger.Info("marshalling user receiver relationship struct")
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

func (rr *RelationshipRepository) GetRelationship(userID string, receiverID string) (*relationship.Relationship, error) {
	rr.logger.Info("getting user receiver relationship from db", zap.String(log.UserIDLogKey, userID), zap.String(log.ReceiverIDLogKey, receiverID))

	result, err := rr.Client.GetItem(rr.Ctx, &dynamodb.GetItemInput{
		TableName: &rr.TableName,
		Key: map[string]types.AttributeValue{
			"user_id":     &types.AttributeValueMemberS{Value: userID},
			"receiver_id": &types.AttributeValueMemberS{Value: receiverID},
		},
	})
	if err != nil {
		return nil, err
	}

	var r relationship.Relationship
	err = attributevalue.UnmarshalMap(result.Item, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (rr *RelationshipRepository) GetRelationshipsByUser(userID string) ([]relationship.Relationship, error) {
	rr.logger.Info("getting user receiver relationships from db", zap.String(log.UserIDLogKey, userID))

	keyCondition := "user_id = :uid"
	expressionAttributeValues := map[string]types.AttributeValue{
		":uid": &types.AttributeValueMemberS{Value: userID},
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(rr.TableName),
		KeyConditionExpression:    aws.String(keyCondition),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	result, err := rr.Client.Query(rr.Ctx, queryInput)
	if err != nil {
		return nil, err
	}

	var relationshipsList []relationship.Relationship
	err = attributevalue.UnmarshalListOfMaps(result.Items, &relationshipsList)
	if err != nil {
		rr.logger.Error("error unmarshalling relationships list", zap.Error(err))
		return nil, err
	}

	return relationshipsList, nil
}

func (rr *RelationshipRepository) DeleteRelationship(userID string, receiverID string) error {
	rr.logger.Info("deleting user receiver relationship from db", zap.String(log.UserIDLogKey, userID), zap.String(log.ReceiverIDLogKey, receiverID))

	_, err := rr.Client.DeleteItem(rr.Ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(rr.TableName),
		Key: map[string]types.AttributeValue{
			"user_id":     &types.AttributeValueMemberS{Value: userID},
			"receiver_id": &types.AttributeValueMemberS{Value: receiverID},
		},
	})

	if err != nil {
		return err
	}

	rr.logger.Info("successfully deleted relationship")
	return nil
}

func (r *RelationshipRepository) GetRelationshipsByEmailNotifications() ([]relationship.Relationship, error) {
	r.logger.Info("getting relationships with email notifications enabled")

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		IndexName:              aws.String("email_notifications"),
		KeyConditionExpression: aws.String("email_notifications_gsi_pk = :email_notifications_gsi_pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email_notifications_gsi_pk": &types.AttributeValueMemberN{Value: "1"},
		},
	}

	var relationships []relationship.Relationship

	paginator := dynamodb.NewQueryPaginator(r.Client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(r.Ctx)
		if err != nil {
			r.logger.Error("failed to query relationships by email notification",
				zap.Error(err))
			return nil, err
		}

		var pageRelationships []relationship.Relationship
		err = attributevalue.UnmarshalListOfMaps(page.Items, &pageRelationships)
		if err != nil {
			r.logger.Error("failed to unmarshal relationships",
				zap.Error(err))
			return nil, err
		}

		relationships = append(relationships, pageRelationships...)
	}

	r.logger.Info("successfully retrieved relationships with email notifications",
		zap.Int("count", len(relationships)))

	return relationships, nil
}

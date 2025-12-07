package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/care-giver-app/care-giver-golang-common/pkg/dynamo"
	"github.com/care-giver-app/care-giver-golang-common/pkg/relationship"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAddRelationship(t *testing.T) {
	tests := map[string]struct {
		relationship *relationship.Relationship
		mockDynamo   *dynamo.Mock
		expectError  bool
	}{
		"Happy Path - Event Added": {
			relationship: &relationship.Relationship{
				UserID:             "User#123",
				ReceiverID:         "Receiver#123",
				PrimaryCareGiver:   true,
				EmailNotifications: false,
			},
			mockDynamo: &dynamo.Mock{
				PutOutput: nil,
				Err:       nil,
			},
		},
		"Sad Path - Put Item Error": {
			relationship: &relationship.Relationship{
				UserID:             "Error",
				ReceiverID:         "Receiver#123",
				PrimaryCareGiver:   true,
				EmailNotifications: false,
			},
			mockDynamo: &dynamo.Mock{
				PutOutput: nil,
				Err:       errors.New("An error occured during Put Item"),
			},
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testEventRepo := NewRelationshipRepository(context.Background(), "relationship-table", tc.mockDynamo, zap.NewNop())

			err := testEventRepo.AddRelationship(tc.relationship)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetRelationship(t *testing.T) {
	tests := map[string]struct {
		userID               string
		receiverID           string
		mockDynamo           *dynamo.Mock
		expectedRelationship relationship.Relationship
		expectError          bool
	}{
		"Happy Path": {
			userID:     "User#123",
			receiverID: "Receiver#123",
			mockDynamo: &dynamo.Mock{
				GetOutput: &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"user_id":             &types.AttributeValueMemberS{Value: "User#123"},
						"receiver_id":         &types.AttributeValueMemberS{Value: "Receiver#123"},
						"primary_care_giver":  &types.AttributeValueMemberBOOL{Value: true},
						"email_notifications": &types.AttributeValueMemberBOOL{Value: false},
					},
				},
				Err: nil,
			},
			expectedRelationship: relationship.Relationship{
				UserID:             "User#123",
				ReceiverID:         "Receiver#123",
				PrimaryCareGiver:   true,
				EmailNotifications: false,
			},
		},
		"Sad Path - Get Item Error": {
			userID:     "Error",
			receiverID: "Receiver#123",
			mockDynamo: &dynamo.Mock{
				GetOutput: nil,
				Err:       errors.New("An error occured during Get Item"),
			},
			expectError: true,
		},
		"Sad Path - Unmarshal Error": {
			userID:     "Unmarshal Error",
			receiverID: "Receiver#123",
			mockDynamo: &dynamo.Mock{
				GetOutput: &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"user_id":            &types.AttributeValueMemberS{Value: "Unmarshal Error"},
						"receiver_id":        &types.AttributeValueMemberS{Value: "testFirstName"},
						"primary_care_giver": &types.AttributeValueMemberS{Value: "false"},
					},
				},
				Err: nil,
			},
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testEventRepo := NewRelationshipRepository(context.Background(), "relationship-table", tc.mockDynamo, zap.NewNop())

			r, err := testEventRepo.GetRelationship(tc.userID, tc.receiverID)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, r)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, r)
				assert.Equal(t, tc.expectedRelationship, *r)
			}
		})
	}
}

func TestGetRelationships(t *testing.T) {
	tests := map[string]struct {
		userID        string
		mockDynamo    *dynamo.Mock
		expectedValue []relationship.Relationship
		expectError   bool
	}{
		"Happy Path": {
			userID: "User#123",
			mockDynamo: &dynamo.Mock{
				QueryOutput: &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{
						{
							"user_id":             &types.AttributeValueMemberS{Value: "User#123"},
							"receiver_id":         &types.AttributeValueMemberS{Value: "Receiver#123"},
							"primary_care_giver":  &types.AttributeValueMemberBOOL{Value: true},
							"email_notifications": &types.AttributeValueMemberBOOL{Value: false},
						},
					},
				},
				Err: nil,
			},
			expectedValue: []relationship.Relationship{
				{
					UserID:             "User#123",
					ReceiverID:         "Receiver#123",
					PrimaryCareGiver:   true,
					EmailNotifications: false,
				},
			},
		},
		"Sad Path - Get Item Error": {
			userID: "Error",
			mockDynamo: &dynamo.Mock{
				QueryOutput: nil,
				Err:         errors.New("An error occured during Query"),
			},
			expectError: true,
		},
		"Sad Path - Unmarshal Error": {
			userID: "Unmarshal Error",
			mockDynamo: &dynamo.Mock{
				QueryOutput: &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{
						{
							"user_id": &types.AttributeValueMemberBOOL{Value: false},
						},
					},
				},
				Err: nil,
			},
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testEventRepo := NewRelationshipRepository(context.Background(), "relationship-table", tc.mockDynamo, zap.NewNop())

			r, err := testEventRepo.GetRelationshipsByUser(tc.userID)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, r)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, r)
				assert.Equal(t, tc.expectedValue, r)
			}
		})
	}
}

func TestDeleteRelationship(t *testing.T) {
	tests := map[string]struct {
		userID      string
		receiverID  string
		mockDynamo  *dynamo.Mock
		expectError bool
	}{
		"Happy Path": {
			userID:     "User#123",
			receiverID: "Receiver#123",
			mockDynamo: &dynamo.Mock{
				DeleteOutput: nil,
				Err:          nil,
			},
		},
		"Sad Path": {
			userID:     "Error",
			receiverID: "Receiver#123",
			mockDynamo: &dynamo.Mock{
				DeleteOutput: nil,
				Err:          errors.New("An error occured during Delete"),
			},
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testEventRepo := NewRelationshipRepository(context.Background(), "relationship-table", tc.mockDynamo, zap.NewNop())

			err := testEventRepo.DeleteRelationship(tc.userID, tc.receiverID)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetRelationshipsByEmailNotifications(t *testing.T) {
	tests := map[string]struct {
		mockDynamo    *dynamo.Mock
		expectedValue []relationship.Relationship
		expectError   bool
	}{
		"Happy Path - Single Page": {
			mockDynamo: &dynamo.Mock{
				QueryOutputs: []*dynamodb.QueryOutput{
					{
						Items: []map[string]types.AttributeValue{
							{
								"user_id":             &types.AttributeValueMemberS{Value: "User#123"},
								"receiver_id":         &types.AttributeValueMemberS{Value: "Receiver#123"},
								"primary_care_giver":  &types.AttributeValueMemberBOOL{Value: true},
								"email_notifications": &types.AttributeValueMemberBOOL{Value: true},
							},
						},
					},
				},
				Err: nil,
			},
			expectedValue: []relationship.Relationship{
				{
					UserID:             "User#123",
					ReceiverID:         "Receiver#123",
					PrimaryCareGiver:   true,
					EmailNotifications: true,
				},
			},
		},
		"Happy Path - Multiple Pages": {
			mockDynamo: &dynamo.Mock{
				QueryOutputs: []*dynamodb.QueryOutput{
					{
						Items: []map[string]types.AttributeValue{
							{
								"user_id":             &types.AttributeValueMemberS{Value: "User#123"},
								"receiver_id":         &types.AttributeValueMemberS{Value: "Receiver#123"},
								"primary_care_giver":  &types.AttributeValueMemberBOOL{Value: true},
								"email_notifications": &types.AttributeValueMemberBOOL{Value: true},
							},
						},
						LastEvaluatedKey: map[string]types.AttributeValue{
							"user_id": &types.AttributeValueMemberS{Value: "User#123"},
						},
					},
					{
						Items: []map[string]types.AttributeValue{
							{
								"user_id":             &types.AttributeValueMemberS{Value: "User#456"},
								"receiver_id":         &types.AttributeValueMemberS{Value: "Receiver#456"},
								"primary_care_giver":  &types.AttributeValueMemberBOOL{Value: false},
								"email_notifications": &types.AttributeValueMemberBOOL{Value: true},
							},
						},
					},
				},
				Err: nil,
			},
			expectedValue: []relationship.Relationship{
				{
					UserID:             "User#123",
					ReceiverID:         "Receiver#123",
					PrimaryCareGiver:   true,
					EmailNotifications: true,
				},
				{
					UserID:             "User#456",
					ReceiverID:         "Receiver#456",
					PrimaryCareGiver:   false,
					EmailNotifications: true,
				},
			},
		},
		"Happy Path - No Results": {
			mockDynamo: &dynamo.Mock{
				QueryOutput: &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{},
				},
				Err: nil,
			},
			expectedValue: nil,
		},
		"Sad Path - Get Item Error": {
			mockDynamo: &dynamo.Mock{
				QueryOutput: nil,
				Err:         errors.New("An error occured during Query"),
			},
			expectError: true,
		},
		"Sad Path - Unmarshal Error": {
			mockDynamo: &dynamo.Mock{
				QueryOutput: &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{
						{
							"user_id": &types.AttributeValueMemberBOOL{Value: false},
						},
					},
				},
				Err: nil,
			},
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testEventRepo := NewRelationshipRepository(context.Background(), "relationship-table", tc.mockDynamo, zap.NewNop())

			r, err := testEventRepo.GetRelationshipsByEmailNotifications()
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, r)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, r)
			}
		})
	}
}

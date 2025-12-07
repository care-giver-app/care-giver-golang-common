package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/care-giver-app/care-giver-golang-common/pkg/dynamo"
	"github.com/care-giver-app/care-giver-golang-common/pkg/user"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateUser(t *testing.T) {
	tests := map[string]struct {
		user        user.User
		mockDynamo  *dynamo.Mock
		expectError bool
	}{
		"Happy Path - User Created": {
			user: user.User{
				UserID:    "User#123",
				FirstName: "testName",
				LastName:  "testLastName",
			},
			mockDynamo: &dynamo.Mock{
				PutOutput: &dynamodb.PutItemOutput{},
				Err:       nil,
			},
			expectError: false,
		},
		"Sad Path - Error Putting Item": {
			user: user.User{
				UserID: "Error",
			},
			mockDynamo: &dynamo.Mock{
				PutOutput: &dynamodb.PutItemOutput{},
				Err:       errors.New("An error occured during Put Item"),
			},
			expectError: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testUserRepo := NewUserRespository(context.Background(), "user-table", tc.mockDynamo, zap.NewNop())

			err := testUserRepo.CreateUser(tc.user)

			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := map[string]struct {
		userID       string
		mockDynamo   *dynamo.Mock
		expectedUser user.User
		expectError  bool
	}{
		"Happy Path - Got User": {
			userID: "User#123",
			mockDynamo: &dynamo.Mock{
				GetOutput: &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"user_id":    &types.AttributeValueMemberS{Value: "User#123"},
						"first_name": &types.AttributeValueMemberS{Value: "testFirstName"},
						"last_name":  &types.AttributeValueMemberS{Value: "testLastName"},
					},
				},
				Err: nil,
			},
			expectedUser: user.User{
				UserID:    "User#123",
				FirstName: "testFirstName",
				LastName:  "testLastName",
			},
		},
		"Sad Path - Error Getting Item": {
			userID: "Get Item Error",
			mockDynamo: &dynamo.Mock{
				GetOutput: nil,
				Err:       errors.New("An error occured during Get Item"),
			},
			expectError: true,
		},
		"Sad Path - Error Unmarshalling Item": {
			userID: "Unmarshal Error",
			mockDynamo: &dynamo.Mock{
				GetOutput: &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"user_id": &types.AttributeValueMemberBOOL{Value: false},
					},
				},
				Err: nil,
			},
			expectError: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testUserRepo := NewUserRespository(context.Background(), "user-table", tc.mockDynamo, zap.NewNop())

			user, err := testUserRepo.GetUser(tc.userID)

			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	tests := map[string]struct {
		email        string
		mockDynamo   *dynamo.Mock
		expectedUser user.User
		expectError  bool
	}{
		"Happy Path - Got User": {
			email: "valid@example.com",
			mockDynamo: &dynamo.Mock{
				QueryOutput: &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{
						{
							"user_id":    &types.AttributeValueMemberS{Value: "User#123"},
							"first_name": &types.AttributeValueMemberS{Value: "testFirstName"},
							"last_name":  &types.AttributeValueMemberS{Value: "testLastName"},
						},
					},
				},
				Err: nil,
			},
			expectedUser: user.User{
				UserID:    "User#123",
				FirstName: "testFirstName",
				LastName:  "testLastName",
			},
		},
		"Sad Path - Error Getting Item": {
			email: "dberror@example.com",
			mockDynamo: &dynamo.Mock{
				QueryOutput: nil,
				Err:         errors.New("An error occured during Query"),
			},
			expectError: true,
		},
		"Sad Path - Error Unmarshalling Item": {
			email: "unmarshalerror@example.com",
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
			testUserRepo := NewUserRespository(context.Background(), "user-table", tc.mockDynamo, zap.NewNop())

			user, err := testUserRepo.GetUserByEmail(tc.email)

			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/care-giver-app/care-giver-golang-common/pkg/dynamo"
	"github.com/care-giver-app/care-giver-golang-common/pkg/receiver"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateReceiver(t *testing.T) {
	tests := map[string]struct {
		receiver    receiver.Receiver
		mockDynamo  *dynamo.Mock
		expectError bool
	}{
		"Happy Path - Receiver Created": {
			receiver: receiver.Receiver{
				ReceiverID: "Receiver#123",
				FirstName:  "testName",
				LastName:   "testLastName",
			},
			mockDynamo: &dynamo.Mock{
				PutOutput: &dynamodb.PutItemOutput{},
				Err:       nil,
			},
		},
		"Sad Path - Error Putting Item": {
			receiver: receiver.Receiver{
				ReceiverID: "Error",
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
			testReceiverRepo := NewReceiverRespository(context.Background(), "receiver-table", tc.mockDynamo, zap.NewNop())

			err := testReceiverRepo.CreateReceiver(tc.receiver)

			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetReceiver(t *testing.T) {
	tests := map[string]struct {
		receiverID       string
		mockDynamo       *dynamo.Mock
		expectedReceiver receiver.Receiver
		expectError      bool
	}{
		"Happy Path - Got Receiver": {
			receiverID: "Receiver#123",
			mockDynamo: &dynamo.Mock{
				GetOutput: &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"receiver_id": &types.AttributeValueMemberS{Value: "Receiver#123"},
						"first_name":  &types.AttributeValueMemberS{Value: "testFirstName"},
						"last_name":   &types.AttributeValueMemberS{Value: "testLastName"},
					},
				},
				Err: nil,
			},
			expectedReceiver: receiver.Receiver{
				ReceiverID: "Receiver#123",
				FirstName:  "testFirstName",
				LastName:   "testLastName",
			},
		},
		"Sad Path - Error Getting Item": {
			receiverID: "Get Item Error",
			mockDynamo: &dynamo.Mock{
				GetOutput: nil,
				Err:       errors.New("An error occured during Get Item"),
			},
			expectError: true,
		},
		"Sad Path - Error Unmarshalling Item": {
			receiverID: "Unmarshal Error",
			mockDynamo: &dynamo.Mock{
				GetOutput: &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"receiver_id": &types.AttributeValueMemberS{Value: "Unmarshal Error"},
						"first_name":  &types.AttributeValueMemberS{Value: "testFirstName"},
						"last_name":   &types.AttributeValueMemberBOOL{Value: false},
					},
				},
				Err: nil,
			},
			expectError: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testReceiverRepo := NewReceiverRespository(context.Background(), "receiver-table", tc.mockDynamo, zap.NewNop())

			receiver, err := testReceiverRepo.GetReceiver(tc.receiverID)
			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedReceiver, receiver)
			}
		})
	}
}

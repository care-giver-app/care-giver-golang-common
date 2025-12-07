package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/care-giver-app/care-giver-golang-common/pkg/dynamo"
	"github.com/care-giver-app/care-giver-golang-common/pkg/event"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAddEvent(t *testing.T) {
	tests := map[string]struct {
		entry       *event.Entry
		mockDynamo  *dynamo.Mock
		expectError bool
	}{
		"Happy Path - Event Added": {
			entry: &event.Entry{
				EventID: "Event#123",
			},
			mockDynamo: &dynamo.Mock{
				PutOutput: &dynamodb.PutItemOutput{},
				Err:       nil,
			},
		},
		"Sad Path - Put Item Error": {
			entry: &event.Entry{
				EventID: "Error",
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
			testEventRepo := NewEventRespository(context.Background(), "event-table", tc.mockDynamo, zap.NewNop())

			err := testEventRepo.AddEvent(tc.entry)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetEvents(t *testing.T) {
	tests := map[string]struct {
		rid           string
		mockDynamo    *dynamo.Mock
		expectedValue []event.Entry
		expectError   bool
	}{
		"Happy Path - Got Events": {
			rid: "Receiver#123",
			mockDynamo: &dynamo.Mock{
				QueryOutput: &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{
						{
							"type": &types.AttributeValueMemberS{Value: "Shower"},
						},
						{
							"type": &types.AttributeValueMemberS{Value: "Medication"},
						},
					},
				},
				Err: nil,
			},
			expectedValue: []event.Entry{
				{
					Type: "Shower",
				},
				{
					Type: "Medication",
				},
			},
		},
		"Sad Path - Query Error": {
			rid: "Error",
			mockDynamo: &dynamo.Mock{
				QueryOutput: nil,
				Err:         errors.New("An error occured during Query"),
			},
			expectError: true,
		},
		"Sad Path - Bad Data Error": {
			rid: "BadData",
			mockDynamo: &dynamo.Mock{
				QueryOutput: &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{
						{
							"type": &types.AttributeValueMemberBOOL{Value: false},
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
			testEventRepo := NewEventRespository(context.Background(), "event-table", tc.mockDynamo, zap.NewNop())

			events, err := testEventRepo.GetEvents(tc.rid)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, events)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, events)
			}
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	tests := map[string]struct {
		rid         string
		eid         string
		mockDynamo  *dynamo.Mock
		expectError bool
	}{
		"Happy Path - Event Deleted": {
			rid: "Receiver#123",
			eid: "Event#123",
			mockDynamo: &dynamo.Mock{
				Err: nil,
			},
		},
		"Sad Path - Delete Error": {
			rid:         "Receiver#123",
			eid:         "Error",
			expectError: true,
			mockDynamo: &dynamo.Mock{
				Err: errors.New("error deleting item"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testEventRepo := NewEventRespository(context.Background(), "event-table", tc.mockDynamo, zap.NewNop())

			err := testEventRepo.DeleteEvent(tc.rid, tc.eid)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

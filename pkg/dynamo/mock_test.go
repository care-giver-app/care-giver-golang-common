package dynamo

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestMock_PutItem(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedOutput := &dynamodb.PutItemOutput{}
		mock := &Mock{PutOutput: expectedOutput}

		output, err := mock.PutItem(ctx, &dynamodb.PutItemInput{})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if output != expectedOutput {
			t.Errorf("expected output %v, got %v", expectedOutput, output)
		}
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("put error")
		mock := &Mock{Err: expectedErr}

		_, err := mock.PutItem(ctx, &dynamodb.PutItemInput{})

		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestMock_GetItem(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedOutput := &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "123"},
			},
		}
		mock := &Mock{GetOutput: expectedOutput}

		output, err := mock.GetItem(ctx, &dynamodb.GetItemInput{})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if output != expectedOutput {
			t.Errorf("expected output %v, got %v", expectedOutput, output)
		}
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("get error")
		mock := &Mock{Err: expectedErr}

		_, err := mock.GetItem(ctx, &dynamodb.GetItemInput{})

		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestMock_UpdateItem(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedOutput := &dynamodb.UpdateItemOutput{}
		mock := &Mock{UpdateOutput: expectedOutput}

		output, err := mock.UpdateItem(ctx, &dynamodb.UpdateItemInput{})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if output != expectedOutput {
			t.Errorf("expected output %v, got %v", expectedOutput, output)
		}
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("update error")
		mock := &Mock{Err: expectedErr}

		_, err := mock.UpdateItem(ctx, &dynamodb.UpdateItemInput{})

		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestMock_Query(t *testing.T) {
	ctx := context.Background()

	t.Run("single output", func(t *testing.T) {
		expectedOutput := &dynamodb.QueryOutput{
			Count: 1,
		}
		mock := &Mock{QueryOutput: expectedOutput}

		output, err := mock.Query(ctx, &dynamodb.QueryInput{})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if output != expectedOutput {
			t.Errorf("expected output %v, got %v", expectedOutput, output)
		}
	})

	t.Run("multiple outputs", func(t *testing.T) {
		output1 := &dynamodb.QueryOutput{Count: 1}
		output2 := &dynamodb.QueryOutput{Count: 2}
		output3 := &dynamodb.QueryOutput{Count: 3}

		mock := &Mock{
			QueryOutputs: []*dynamodb.QueryOutput{output1, output2, output3},
		}

		result1, err := mock.Query(ctx, &dynamodb.QueryInput{})
		if err != nil {
			t.Errorf("expected no error on first call, got %v", err)
		}
		if result1 != output1 {
			t.Errorf("expected first output, got %v", result1)
		}

		result2, err := mock.Query(ctx, &dynamodb.QueryInput{})
		if err != nil {
			t.Errorf("expected no error on second call, got %v", err)
		}
		if result2 != output2 {
			t.Errorf("expected second output, got %v", result2)
		}

		result3, err := mock.Query(ctx, &dynamodb.QueryInput{})
		if err != nil {
			t.Errorf("expected no error on third call, got %v", err)
		}
		if result3 != output3 {
			t.Errorf("expected third output, got %v", result3)
		}
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("query error")
		mock := &Mock{Err: expectedErr}

		_, err := mock.Query(ctx, &dynamodb.QueryInput{})

		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestMock_DeleteItem(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedOutput := &dynamodb.DeleteItemOutput{}
		mock := &Mock{DeleteOutput: expectedOutput}

		output, err := mock.DeleteItem(ctx, &dynamodb.DeleteItemInput{})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if output != expectedOutput {
			t.Errorf("expected output %v, got %v", expectedOutput, output)
		}
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("delete error")
		mock := &Mock{Err: expectedErr}

		_, err := mock.DeleteItem(ctx, &dynamodb.DeleteItemInput{})

		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}

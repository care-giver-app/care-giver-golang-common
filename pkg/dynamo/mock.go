package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Mock struct {
	Err          error
	QueryOutput  *dynamodb.QueryOutput
	QueryOutputs []*dynamodb.QueryOutput
	QueryCallNum int
	PutOutput    *dynamodb.PutItemOutput
	GetOutput    *dynamodb.GetItemOutput
	UpdateOutput *dynamodb.UpdateItemOutput
	DeleteOutput *dynamodb.DeleteItemOutput
}

func (m *Mock) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.PutOutput, m.Err
}

func (m *Mock) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return m.GetOutput, m.Err
}

func (m *Mock) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	return m.UpdateOutput, m.Err
}

func (m *Mock) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	if len(m.QueryOutputs) > 0 {
		output := m.QueryOutputs[m.QueryCallNum]
		m.QueryCallNum++
		return output, m.Err
	}

	return m.QueryOutput, m.Err
}

func (m *Mock) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return m.DeleteOutput, m.Err
}

package awsconfig

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

const (
	USEastTwoRegion = "us-east-2"
	localEnv        = "local"
)

func GetAWSConfig(ctx context.Context, env string) (aws.Config, error) {
	if env == localEnv {
		return getLocalAWSConfig(ctx)
	}
	return getAWSConfig(ctx)
}

func getAWSConfig(ctx context.Context) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(USEastTwoRegion),
	)

	return cfg, err
}

func getLocalAWSConfig(ctx context.Context) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(USEastTwoRegion),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)

	return cfg, err
}

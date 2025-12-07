package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateClient(t *testing.T) {
	tests := map[string]struct {
		env               string
		awsConfig         aws.Config
		logger            *zap.Logger
		expectedEndpoint  string
		expectNilEndpoint bool
	}{
		"Happy Path - Local Endpoint": {
			env:              localEnv,
			awsConfig:        aws.Config{},
			logger:           zap.Must(zap.NewProduction()),
			expectedEndpoint: localDockerEndpoint,
		},
		"Happy Path - Non Local Endpoint": {
			env:               "dev",
			awsConfig:         aws.Config{},
			logger:            zap.Must(zap.NewProduction()),
			expectNilEndpoint: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := CreateClient(tc.env, tc.awsConfig, tc.logger)

			if tc.expectNilEndpoint {
				assert.Nil(t, client.Options().BaseEndpoint)
			} else {
				assert.Equal(t, tc.expectedEndpoint, *client.Options().BaseEndpoint)
			}
		})
	}
}

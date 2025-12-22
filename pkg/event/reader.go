package event

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	configDirectory = "types/"
	jsonSuffix      = ".json"
)

//go:embed types/*.json
var configs embed.FS

type EventConfig struct {
	Type string     `json:"type"`
	Data DataConfig `json:"data,omitempty"`
}

type DataConfig struct {
	Name string `json:"name"`
	Unit string `json:"unit"`
}

func readEventConfig(eventType string) (*EventConfig, error) {
	fileName := prepareFileName(eventType)

	configBytes, err := configs.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var eventConfig EventConfig
	err = json.Unmarshal(configBytes, &eventConfig)
	if err != nil {
		return nil, err
	}

	return &eventConfig, nil
}

func prepareFileName(eventType string) string {
	preparedType := strings.ReplaceAll(strings.ToLower(eventType), " ", "_")
	return fmt.Sprintf("%s%s%s", configDirectory, preparedType, jsonSuffix)
}

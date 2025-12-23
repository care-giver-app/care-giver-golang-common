package event

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	configDirectory = "types"
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
	name := prepareFileName(eventType)

	content, err := configs.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var eventConfig EventConfig
	err = json.Unmarshal(content, &eventConfig)
	if err != nil {
		return nil, err
	}

	return &eventConfig, nil
}

func prepareFileName(eventType string) string {
	preparedType := strings.ReplaceAll(strings.ToLower(eventType), " ", "_")
	return fmt.Sprintf("%s/%s%s", configDirectory, preparedType, jsonSuffix)
}

func readConfigs() ([]EventConfig, error) {
	entries, err := configs.ReadDir(configDirectory)
	if err != nil {
		return nil, err
	}

	var eventConfigs []EventConfig
	for _, entry := range entries {
		name := entry.Name()
		content, err := configs.ReadFile(fmt.Sprintf("%s/%s", configDirectory, name))
		if err != nil {
			return nil, err
		}

		var eventConfig EventConfig
		err = json.Unmarshal(content, &eventConfig)
		if err != nil {
			return nil, err
		}

		eventConfigs = append(eventConfigs, eventConfig)
	}

	return eventConfigs, nil
}

package event

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEntry(t *testing.T) {
	tests := map[string]struct {
		eventType     string
		startTime     string
		endTime       string
		data          []DataPoint
		note          string
		expectedEntry Entry
		expectErr     bool
	}{
		"Happy Path": {
			eventType: "Shower",
			startTime: "2025-06-30T17:25:00-05:00",
			endTime:   "2025-06-30T17:55:00-05:00",
			expectedEntry: Entry{
				Type:      "Shower",
				StartTime: "2025-06-30T17:25:00-05:00",
				EndTime:   "2025-06-30T17:55:00-05:00",
			},
		},
		"Happy Path - With End Time": {
			eventType: "Medication",
			startTime: "2025-06-30T17:25:00-05:00",
			endTime:   "2026-06-30T17:55:00-05:00",
			expectedEntry: Entry{
				Type:      "Medication",
				StartTime: "2025-06-30T17:25:00-05:00",
				EndTime:   "2026-06-30T17:55:00-05:00",
			},
		},
		"Happy Path - With Data": {
			eventType: "Weight",
			data: []DataPoint{
				{
					Name:  "Weight",
					Value: 120.3,
				},
			},
			startTime: "2025-06-30T17:25:00-05:00",
			endTime:   "2025-06-30T17:55:00-05:00",
			expectedEntry: Entry{
				Type:      "Weight",
				StartTime: "2025-06-30T17:25:00-05:00",
				EndTime:   "2025-06-30T17:55:00-05:00",
				Data: []DataPoint{
					{
						Name:  "Weight",
						Value: 120.3,
					},
				},
			},
		},
		"Happy Path - With Note": {
			eventType: "Weight",
			startTime: "2025-06-30T17:25:00-05:00",
			endTime:   "2025-06-30T17:55:00-05:00",
			note:      "some note",
			expectedEntry: Entry{
				Type:      "Weight",
				StartTime: "2025-06-30T17:25:00-05:00",
				EndTime:   "2025-06-30T17:55:00-05:00",
				Note:      "some note",
			},
		},
		"Sad Path - Bad Event Type": {
			eventType: "BadEventType",
			expectErr: true,
		},
	}

	testRID := "Receiver#123"
	testUID := "User#123"

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.expectedEntry.ReceiverID = testRID
			tc.expectedEntry.UserID = testUID

			opts := []EntryOption{}

			if len(tc.data) > 0 {
				opts = append(opts, WithData(tc.data))
				tc.expectedEntry.Data = tc.data
			}

			if tc.note != "" {
				opts = append(opts, WithNote(tc.note))
			}

			entry, err := NewEntry(testRID, testUID, tc.eventType, tc.startTime, tc.endTime, opts...)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, entry)
			} else {
				tc.expectedEntry.EventID = entry.EventID
				tc.expectedEntry.StartTime = entry.StartTime
				tc.expectedEntry.EndTime = entry.EndTime
				assert.Equal(t, tc.expectedEntry, *entry)
			}
		})
	}

}

func TestGetAllConfigs(t *testing.T) {
	tests := map[string]struct {
		expectedConfigs []EventConfig
		expectErr       bool
	}{
		"Happy Path": {
			expectedConfigs: []EventConfig{
				{
					Type: "Shower",
					Icon: "assets/shower-icon.png",
					Color: ColorConfig{
						Primary:   "#3498DB",
						Secondary: "#D6EAF8",
					},
				},
				{
					Type: "Weight",
					Icon: "assets/weight-icon.png",
					Color: ColorConfig{
						Primary:   "#27AE60",
						Secondary: "#D4EFDF",
					},
					Data: &DataConfig{
						Name: "Weight",
						Unit: "Lbs",
					},
					Graph: &GraphConfig{
						Type:  "line",
						Title: "Weight By Time",
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			configs, err := GetAllConfigs()

			assert.Nil(t, err)
			assert.NotEmpty(t, configs)

			for _, expectedConfig := range tc.expectedConfigs {
				found := false
				for _, actualConfig := range configs {
					if actualConfig.Type == expectedConfig.Type {
						found = true
						break
					}
				}
				if !found {
					assert.Fail(t, fmt.Sprintf("%s not found in configs", expectedConfig.Type))
				}
			}
		})
	}
}

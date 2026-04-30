package event

import (
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
					Type: "Doctor Appointment",
					Icon: "assets/appointment-icon.svg",
					Color: ColorConfig{
						Primary:   "#1565C0",
						Secondary: "#BBDEFB",
					},
					Fields: []FieldConfig{
						{Name: "Doctor", Label: "Doctor / Provider", InputType: "text", Required: true, Placeholder: "e.g. Dr. Smith"},
						{Name: "Specialty", Label: "Specialty", InputType: "text", Required: false, Placeholder: "e.g. Cardiology"},
						{Name: "Location", Label: "Clinic / Location", InputType: "text", Required: false, Placeholder: "e.g. Cleveland Clinic"},
						{Name: "Reason", Label: "Reason for Visit", InputType: "text", Required: false, Placeholder: "e.g. Annual checkup"},
						{Name: "Outcome", Label: "Outcome / Summary", InputType: "textarea", Required: false, Placeholder: "What happened at the appointment?"},
						{Name: "FollowUp", Label: "Follow-up Date", InputType: "date", Required: false, Placeholder: ""},
					},
				},
				{
					Type: "Shower",
					Icon: "assets/shower-icon.svg",
					Color: ColorConfig{
						Primary:   "#3498DB",
						Secondary: "#D6EAF8",
					},
				},
				{
					Type: "Weight",
					Icon: "assets/weight-icon.svg",
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

	t.Run("alert-mode monitor - Urination", func(t *testing.T) {
		c, ok := byType["Urination"]
		assert.True(t, ok)
		assert.True(t, c.HasQuickAdd)
		assert.NotNil(t, c.Monitor)
		assert.NotNil(t, c.Monitor.AlertThresholds)
		assert.Equal(t, 4, c.Monitor.AlertThresholds.Yellow)
		assert.Equal(t, 8, c.Monitor.AlertThresholds.Red)
		assert.Equal(t, 12, c.Monitor.AlertThresholds.Critical)
		assert.False(t, c.Monitor.ShowLastValue)
	})

	t.Run("alert-mode monitor - Medication", func(t *testing.T) {
		c, ok := byType["Medication"]
		assert.True(t, ok)
		assert.True(t, c.HasQuickAdd)
		assert.NotNil(t, c.Monitor)
		assert.NotNil(t, c.Monitor.AlertThresholds)
		assert.Equal(t, 6, c.Monitor.AlertThresholds.Yellow)
		assert.Equal(t, 12, c.Monitor.AlertThresholds.Red)
		assert.Equal(t, 18, c.Monitor.AlertThresholds.Critical)
	})

	t.Run("alert-mode monitor - Bowel Movement", func(t *testing.T) {
		c, ok := byType["Bowel Movement"]
		assert.True(t, ok)
		assert.True(t, c.HasQuickAdd)
		assert.NotNil(t, c.Monitor)
		assert.NotNil(t, c.Monitor.AlertThresholds)
		assert.Equal(t, 24, c.Monitor.AlertThresholds.Yellow)
		assert.Equal(t, 48, c.Monitor.AlertThresholds.Red)
		assert.Equal(t, 72, c.Monitor.AlertThresholds.Critical)
	})

	t.Run("alert-mode monitor - Shower", func(t *testing.T) {
		c, ok := byType["Shower"]
		assert.True(t, ok)
		assert.True(t, c.HasQuickAdd)
		assert.NotNil(t, c.Monitor)
		assert.NotNil(t, c.Monitor.AlertThresholds)
		assert.Equal(t, 36, c.Monitor.AlertThresholds.Yellow)
		assert.Equal(t, 60, c.Monitor.AlertThresholds.Red)
		assert.Equal(t, 84, c.Monitor.AlertThresholds.Critical)
	})

	t.Run("last-value monitor - Walk", func(t *testing.T) {
		c, ok := byType["Walk"]
		assert.True(t, ok)
		assert.True(t, c.HasQuickAdd)
		assert.NotNil(t, c.Monitor)
		assert.True(t, c.Monitor.ShowLastValue)
		assert.Nil(t, c.Monitor.AlertThresholds)
	})

	t.Run("last-value monitor - Weight", func(t *testing.T) {
		c, ok := byType["Weight"]
		assert.True(t, ok)
		assert.True(t, c.HasQuickAdd)
		assert.NotNil(t, c.Monitor)
		assert.True(t, c.Monitor.ShowLastValue)
		assert.Nil(t, c.Monitor.AlertThresholds)
	})

	t.Run("upcoming - Doctor Appointment", func(t *testing.T) {
		c, ok := byType["Doctor Appointment"]
		assert.True(t, ok)
		assert.False(t, c.HasQuickAdd)
		assert.Nil(t, c.Monitor)
		assert.NotNil(t, c.Upcoming)
		assert.True(t, c.Upcoming.Show)
		assert.Equal(t, 30, c.Upcoming.LookAheadDays)
	})
}

package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	DBPrefix = "Event"
	ParamID  = "eventId"
)

type Entry struct {
	EventID     string      `json:"eventId" dynamodbav:"event_id"`
	ReceiverID  string      `json:"receiverId" dynamodbav:"receiver_id"`
	UserID      string      `json:"userId" dynamodbav:"user_id"`
	StartTime   string      `json:"startTime" dynamodbav:"start_time"`
	EndTime     string      `json:"endTime" dynamodbav:"end_time"`
	Type        string      `json:"type" dynamodbav:"type"`
	IsTrackable bool        `json:"isTrackable" dynamodbav:"is_trackable"`
	Data        []DataPoint `json:"data,omitempty" dynamodbav:"data,omitempty"`
	Note        string      `json:"note,omitempty" dynamodbav:"note,omitempty"`
}

type DataPoint struct {
	Name  string `json:"name" dynamodbav:"name"`
	Value any    `json:"value" dynamodbav:"value"`
}

type EntryOption func(*Entry)

func WithEndTime(endTime string) EntryOption {
	return func(e *Entry) {
		e.EndTime = endTime
	}
}

func WithData(data []DataPoint) EntryOption {
	return func(e *Entry) {
		e.Data = data
	}
}

func WithNote(note string) EntryOption {
	return func(e *Entry) {
		e.Note = note
	}
}

func NewEntry(receiverID, userID, eventType, startTime string, opts ...EntryOption) (*Entry, error) {
	eventConfig, err := readEventConfig(eventType)
	if err != nil {
		return nil, err
	}

	st, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, err
	}

	endTime := st.Add(30 * time.Minute)

	e := &Entry{
		EventID:    fmt.Sprintf("%s#%s", DBPrefix, uuid.New().String()),
		ReceiverID: receiverID,
		UserID:     userID,
		StartTime:  startTime,
		EndTime:    endTime.Format(time.RFC3339),
		Type:       eventConfig.Type,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e, nil
}

func GetAllConfigs() ([]EventConfig, error) {
	return readConfigs()
}

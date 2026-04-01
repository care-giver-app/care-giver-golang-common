package event

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	DBPrefix = "Event"
	ParamID  = "eventId"
)

type Entry struct {
	EventID    string      `json:"eventId" dynamodbav:"event_id"`
	ReceiverID string      `json:"receiverId" dynamodbav:"receiver_id"`
	UserID     string      `json:"userId" dynamodbav:"user_id"`
	StartTime  string      `json:"startTime" dynamodbav:"start_time"`
	EndTime    string      `json:"endTime" dynamodbav:"end_time"`
	Type       string      `json:"type" dynamodbav:"type"`
	Data       []DataPoint `json:"data,omitempty" dynamodbav:"data,omitempty"`
	Note       string      `json:"note,omitempty" dynamodbav:"note,omitempty"`
}

type DataPoint struct {
	Name  string `json:"name" dynamodbav:"name"`
	Value any    `json:"value" dynamodbav:"value"`
}

type EntryOption func(*Entry)

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

func NewEntry(receiverID, userID, eventType, startTime, endTime string, opts ...EntryOption) (*Entry, error) {
	eventConfig, err := readEventConfig(eventType)
	if err != nil {
		return nil, err
	}

	e := &Entry{
		EventID:    fmt.Sprintf("%s#%s", DBPrefix, uuid.New().String()),
		ReceiverID: receiverID,
		UserID:     userID,
		StartTime:  startTime,
		EndTime:    endTime,
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

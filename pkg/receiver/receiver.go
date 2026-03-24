package receiver

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	DBPrefix = "Receiver"
	ParamID  = "receiverId"
)

type Receiver struct {
	ReceiverID     string   `json:"receiverId" dynamodbav:"receiver_id"`
	FirstName      string   `json:"firstName" dynamodbav:"first_name"`
	LastName       string   `json:"lastName" dynamodbav:"last_name"`
	EventsTracking []string `json:"eventsTracking" dynamodbav:"events_tracking"`
}

type ReceiverOption func(*Receiver)

func WithEventsTracking(events []string) ReceiverOption {
	return func(r *Receiver) {
		r.EventsTracking = events
	}
}

func NewReceiver(firstName string, lastName string, opts ...ReceiverOption) *Receiver {
	r := &Receiver{
		ReceiverID: fmt.Sprintf("%s#%s", DBPrefix, uuid.New()),
		FirstName:  firstName,
		LastName:   lastName,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

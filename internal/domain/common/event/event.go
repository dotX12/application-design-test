package event

import (
	"github.com/google/uuid"
	"time"
)

type Contract interface {
	String() string
}

type Event struct {
	EventID uuid.UUID
	Time    time.Time
}

func NewEvent() Event {
	return Event{
		EventID: uuid.New(),
		Time:    time.Now(),
	}
}

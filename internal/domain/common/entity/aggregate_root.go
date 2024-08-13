package entity

import (
	"applicationDesignTest/internal/domain/common/event"
	"sync"
)

type AggregateRoot struct {
	events []event.Contract
	mu     sync.RWMutex
}

func (a *AggregateRoot) RecordEvent(e event.Contract) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = append(a.events, e)
}

func (a *AggregateRoot) GetEvents() []event.Contract {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.events
}

func (a *AggregateRoot) ClearEvents() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = []event.Contract{}
}

func (a *AggregateRoot) PullEvents() []event.Contract {
	a.mu.Lock()
	defer a.mu.Unlock()
	events := a.events
	a.events = []event.Contract{}
	return events
}

package event

import (
	"applicationDesignTest/internal/domain/common/event"
	"time"
)

type OrderCreated struct {
	event.Event
	OrderID string
	Email   string
	From    time.Time
	To      time.Time
}

func NewOrderCreated(orderID, email string, from, to time.Time) OrderCreated {
	return OrderCreated{
		Event:   event.NewEvent(),
		OrderID: orderID,
		Email:   email,
		From:    from,
		To:      to,
	}
}

func (o OrderCreated) String() string {
	return "booking.OrderCreated"
}

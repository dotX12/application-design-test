package event

import (
	"applicationDesignTest/internal/domain/event"
	"applicationDesignTest/internal/observability"
	"context"
)

type NotificationOrderCreatedHandler struct {
	observer *observability.Observability
}

func NewNotificationOrderCreatedHandler(observer *observability.Observability) NotificationOrderCreatedHandler {
	return NotificationOrderCreatedHandler{
		observer: observer,
	}
}

func (h NotificationOrderCreatedHandler) Handle(ctx context.Context, notification event.OrderCreated) error {
	h.observer.Logger.
		Trace().
		Ctx(ctx).
		Str("event_id", notification.Event.EventID.String()).
		Str("topic", notification.String()).
		Str("created_at", notification.Event.Time.String()).
		Str("order_id", notification.OrderID).
		Str("email", notification.Email).
		Time("from", notification.From).
		Time("to", notification.To).
		Msg("We send a email about the hotel reservation...")
	return nil
}

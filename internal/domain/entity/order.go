package entity

import (
	"applicationDesignTest/internal/domain/common/entity"
	"applicationDesignTest/internal/domain/event"
	"applicationDesignTest/internal/domain/exception"
	"applicationDesignTest/internal/domain/vo"
	"fmt"
	"time"
)

type Order struct {
	entity.AggregateRoot
	ID        vo.OrderID
	HotelID   vo.HotelID
	RoomID    vo.RoomID
	UserEmail vo.Email
	From      time.Time
	To        time.Time
}

func NewOrder(hotelID vo.HotelID, roomID vo.RoomID, userEmail vo.Email, from, to time.Time) (*Order, error) {
	if to.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, fmt.Errorf("date cannot be in the past")
	}

	if from.After(to) {
		return nil, exception.ErrInvalidDateRange
	}
	order := &Order{
		ID:        vo.NewOrderID(),
		HotelID:   hotelID,
		RoomID:    roomID,
		UserEmail: userEmail,
		From:      from,
		To:        to,
	}
	order.RecordEvent(
		event.NewOrderCreated(
			order.ID.String(),
			order.UserEmail.String(),
			order.From,
			order.To,
		),
	)
	return order, nil
}

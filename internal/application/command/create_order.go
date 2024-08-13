package command

import (
	"applicationDesignTest/internal/application/common/exception"
	"applicationDesignTest/internal/domain/entity"
	"applicationDesignTest/internal/domain/vo"
	"applicationDesignTest/internal/observability"
	"applicationDesignTest/internal/port/storage"
	"applicationDesignTest/pkg/mediator"
	"applicationDesignTest/pkg/timeutils"
	"context"
	"fmt"
	"time"
)

type AddOrderCommand struct {
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

type AddOrderCommandHandler struct {
	observer                *observability.Observability
	orderStorage            storage.OrderStorage
	roomAvailabilityStorage storage.RoomAvailabilityStorage
}

func NewAddOrderCommandHandler(
	observer *observability.Observability,
	orderStorage storage.OrderStorage,
	roomAvailabilityStorage storage.RoomAvailabilityStorage,
) *AddOrderCommandHandler {
	return &AddOrderCommandHandler{
		observer:                observer,
		orderStorage:            orderStorage,
		roomAvailabilityStorage: roomAvailabilityStorage,
	}
}

func (h AddOrderCommandHandler) Handle(ctx context.Context, cmd *AddOrderCommand) (*entity.Order, error) {
	h.observer.Logger.Trace().Ctx(ctx).Msg("AddOrderCommandHandler.Handle")
	hotelIDVO := vo.NewHotelIDFromString(cmd.HotelID)

	roomIDVO := vo.NewRoomIDFromString(cmd.RoomID)

	emailVO, err := vo.NewEmailFromString(cmd.UserEmail)
	if err != nil {
		return nil, err
	}

	daysToBook, err := timeutils.DaysBetween(cmd.From, cmd.To)
	if err != nil {
		return nil, err
	}

	for _, day := range daysToBook {
		id, err := vo.NewRoomAvailabilityID(hotelIDVO.String(), roomIDVO.String(), day)
		if err != nil {
			h.observer.Logger.Error().Ctx(ctx).
				Str("hotel_id", hotelIDVO.String()).
				Str("room_id", roomIDVO.String()).
				Time("date", day).
				Msg("Error creating room availability id")
			return nil, err
		}

		availability, err := h.roomAvailabilityStorage.FindByID(ctx, *id)
		if err != nil {
			h.observer.Logger.Error().Ctx(ctx).
				Str("hotel_id", hotelIDVO.String()).
				Str("room_id", roomIDVO.String()).
				Time("date", day).
				Msg("Error finding room availability")
			return nil, fmt.Errorf("%w: %w", err, exception.ErrFindRoomAvailability)
		}
		if !availability.IsAvailable() {
			h.observer.Logger.Error().Ctx(ctx).
				Str("hotel_id", hotelIDVO.String()).
				Str("room_id", roomIDVO.String()).
				Time("date", day).
				Msg("Room not available")
			return nil, exception.ErrHotelRoomNotAvailable
		}
	}

	for _, day := range daysToBook {
		id, err := vo.NewRoomAvailabilityID(hotelIDVO.String(), roomIDVO.String(), day)
		if err != nil {
			h.observer.Logger.Error().Ctx(ctx).
				Str("hotel_id", hotelIDVO.String()).
				Str("room_id", roomIDVO.String()).
				Time("date", day).
				Msg("Error creating room availability id")
			return nil, err
		}

		availability, err := h.roomAvailabilityStorage.FindByID(ctx, *id)
		if err != nil {
			h.observer.Logger.Error().Ctx(ctx).
				Str("hotel_id", hotelIDVO.String()).
				Str("room_id", roomIDVO.String()).
				Time("date", day).
				Msg("Error finding room availability")
			return nil, fmt.Errorf("%w: %w", err, exception.ErrFindRoomAvailability)
		}

		err = availability.Reserve()
		if err != nil {
			return nil, err
		}
		err = h.roomAvailabilityStorage.Update(ctx, availability)
		if err != nil {
			h.observer.Logger.
				Error().Err(err).Ctx(ctx).
				Str("hotel_id", hotelIDVO.String()).
				Str("room_id", roomIDVO.String()).
				Time("date", day).
				Msg("Error updating room availability")
			return nil, fmt.Errorf("%w: %w", err, exception.ErrSaveRoomAvailability)
		}
	}

	orderEntity, err := entity.NewOrder(hotelIDVO, roomIDVO, emailVO, cmd.From, cmd.To)
	if err != nil {
		h.observer.Logger.Error().Err(err).Ctx(ctx).Msg("Error creating order")
		return nil, err
	}

	err = h.orderStorage.Create(ctx, orderEntity)
	if err != nil {
		h.observer.Logger.
			Error().
			Err(err).
			Ctx(ctx).
			Str("hotel_id", hotelIDVO.String()).
			Str("room_id", roomIDVO.String()).
			Str("email", emailVO.String()).
			Time("from", cmd.From).
			Time("to", cmd.To).
			Msg("Error saving order")
		return nil, fmt.Errorf("%w: %w", err, exception.ErrSaveOrder)
	}

	_ = mediator.PublishAsync(ctx, orderEntity.PullEvents()...)

	return orderEntity, nil
}

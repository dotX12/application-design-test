package command

import (
	"applicationDesignTest/internal/adapter/storage/repository"
	"applicationDesignTest/internal/domain/entity"
	"applicationDesignTest/internal/domain/vo"
	"applicationDesignTest/internal/observability"
	"context"
	"time"
)

type AddRoomAvailabilityCommand struct {
	RoomID  string
	HotelID string
	Date    time.Time
	Quota   int
}

type RoomAvailabilityCommandHandler struct {
	observer                *observability.Observability
	roomAvailabilityStorage *repository.ImMemoryRoomAvailabilityStorage
}

func NewRoomAvailabilityCommandHandler(
	observer *observability.Observability,
	roomAvailabilityStorage *repository.ImMemoryRoomAvailabilityStorage,
) *RoomAvailabilityCommandHandler {
	return &RoomAvailabilityCommandHandler{
		observer:                observer,
		roomAvailabilityStorage: roomAvailabilityStorage,
	}
}

func (h RoomAvailabilityCommandHandler) Handle(ctx context.Context, cmd *AddRoomAvailabilityCommand) (*vo.RoomAvailabilityID, error) {
	h.observer.Logger.Info().Ctx(ctx).Msg("RoomAvailabilityCommandHandler.Handle")

	roomVO := vo.NewRoomIDFromString(cmd.RoomID)
	hotelVO := vo.NewHotelIDFromString(cmd.HotelID)
	quota, err := vo.NewQuotaFromInt(cmd.Quota)
	if err != nil {
		return nil, err
	}
	room, err := entity.NewRoomAvailability(hotelVO, roomVO, cmd.Date, *quota)
	if err != nil {
		return nil, err
	}

	err = h.roomAvailabilityStorage.Save(ctx, room)
	if err != nil {
		return nil, err
	}
	return &room.ID, nil
}

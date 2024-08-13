package storage

import (
	"applicationDesignTest/internal/domain/entity"
	"applicationDesignTest/internal/domain/vo"
	"context"
)

type RoomAvailabilityStorage interface {
	IsAvailable(ctx context.Context, id vo.RoomAvailabilityID) (bool, error)
	FindByID(ctx context.Context, id vo.RoomAvailabilityID) (*entity.RoomAvailability, error)
	Save(ctx context.Context, roomAvailability *entity.RoomAvailability) error
	Update(ctx context.Context, roomAvailability *entity.RoomAvailability) error
}

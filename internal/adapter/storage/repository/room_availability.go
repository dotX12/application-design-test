package repository

import (
	"applicationDesignTest/internal/application/common/exception"
	"applicationDesignTest/internal/domain/entity"
	"applicationDesignTest/internal/domain/vo"
	"applicationDesignTest/pkg/logging"
	"applicationDesignTest/pkg/timeutils"
	"context"
	"sync"
	"time"
)

type ImMemoryRoomAvailabilityStorage struct {
	roomAvailabilities map[vo.RoomAvailabilityID]*entity.RoomAvailability
	mu                 sync.RWMutex
}

func (r *ImMemoryRoomAvailabilityStorage) IsAvailable(ctx context.Context, id vo.RoomAvailabilityID) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	roomAvailability, ok := r.roomAvailabilities[id]
	if !ok {
		return false, exception.ErrRoomAvailabilityNotFound
	}
	return roomAvailability.IsAvailable(), nil
}

func (r *ImMemoryRoomAvailabilityStorage) FindByID(ctx context.Context, id vo.RoomAvailabilityID) (*entity.RoomAvailability, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	roomAvailability, ok := r.roomAvailabilities[id]
	if !ok {
		return nil, exception.ErrRoomAvailabilityNotFound
	}
	return roomAvailability, nil
}

func (r *ImMemoryRoomAvailabilityStorage) Save(ctx context.Context, roomAvailability *entity.RoomAvailability) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.roomAvailabilities[roomAvailability.ID]; exists {
		return exception.ErrRoomAvailabilityAlreadyExists
	}
	r.roomAvailabilities[roomAvailability.ID] = roomAvailability
	return nil
}

func (r *ImMemoryRoomAvailabilityStorage) Update(ctx context.Context, roomAvailability *entity.RoomAvailability) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.roomAvailabilities[roomAvailability.ID]; !exists {
		return exception.ErrRoomAvailabilityNotFound
	}
	r.roomAvailabilities[roomAvailability.ID] = roomAvailability
	return nil
}

func NewImMemoryRoomAvailabilityStorage() *ImMemoryRoomAvailabilityStorage {
	storage := &ImMemoryRoomAvailabilityStorage{
		roomAvailabilities: make(map[vo.RoomAvailabilityID]*entity.RoomAvailability),
	}

	initialData := []struct {
		hotelID string
		roomID  string
		date    time.Time
		quota   int
	}{
		{"reddison", "lux", timeutils.DateToTime(2024, 9, 12), 1},
		{"reddison", "lux", timeutils.DateToTime(2024, 9, 13), 1},
		{"reddison", "lux", timeutils.DateToTime(2024, 9, 14), 1},
		{"reddison", "lux", timeutils.DateToTime(2024, 9, 15), 1},
		{"reddison", "lux", timeutils.DateToTime(2024, 9, 16), 0},
	}

	for _, data := range initialData {
		hotelID := vo.NewHotelIDFromString(data.hotelID)
		roomID := vo.NewRoomIDFromString(data.roomID)
		quota, err := vo.NewQuotaFromInt(data.quota)
		if err != nil {
			logging.DefaultLogger.Info().
				Msgf(
					"Error initializing Quota for hotel %s, room %s on date %s: %v",
					data.hotelID,
					data.roomID,
					data.date,
					err,
				)
			continue
		}

		roomAvailability, err := entity.NewRoomAvailability(hotelID, roomID, data.date, *quota)
		if err != nil {
			logging.DefaultLogger.Error().
				Err(err).
				Msg("Error initializing RoomAvailability")
			continue
		}
		storage.roomAvailabilities[roomAvailability.ID] = roomAvailability
	}

	return storage
}

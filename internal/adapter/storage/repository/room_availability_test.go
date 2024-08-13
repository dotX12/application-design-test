package repository_test

import (
	"applicationDesignTest/internal/adapter/storage/repository"
	"applicationDesignTest/internal/application/common/exception"
	"applicationDesignTest/internal/domain/entity"
	"applicationDesignTest/internal/domain/vo"
	"applicationDesignTest/pkg/timeutils"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImMemoryRoomAvailabilityStorage_CheckAvailability(t *testing.T) {
	storage := repository.NewImMemoryRoomAvailabilityStorage()

	id, err := vo.NewRoomAvailabilityID("reddison", "lux", timeutils.DateToTime(2024, 9, 12))
	assert.NoError(t, err)

	isAvailable, err := storage.IsAvailable(context.Background(), *id)
	assert.NoError(t, err)
	assert.True(t, isAvailable)

	id, err = vo.NewRoomAvailabilityID("reddison", "lux", timeutils.DateToTime(2024, 9, 15))
	assert.NoError(t, err)
	isAvailable, err = storage.IsAvailable(context.Background(), *id)
	assert.NoError(t, err)
	assert.True(t, isAvailable)

	id, err = vo.NewRoomAvailabilityID("reddison", "lux", timeutils.DateToTime(2024, 9, 16))
	assert.NoError(t, err)
	isAvailable, err = storage.IsAvailable(context.Background(), *id)
	assert.NoError(t, err)
	assert.False(t, isAvailable)

	id, err = vo.NewRoomAvailabilityID("reddison", "lux", timeutils.DateToTime(2024, 9, 22))
	assert.NoError(t, err)
	isAvailable, err = storage.IsAvailable(context.Background(), *id)
	assert.Error(t, err)
	assert.False(t, isAvailable)
}

func TestImMemoryRoomAvailabilityStorage_FindByID(t *testing.T) {
	storage := repository.NewImMemoryRoomAvailabilityStorage()

	id, err := vo.NewRoomAvailabilityID("reddison", "lux", timeutils.DateToTime(2024, 9, 12))
	assert.NoError(t, err)
	roomAvailability, err := storage.FindByID(context.Background(), *id)
	assert.NoError(t, err)
	assert.NotNil(t, roomAvailability)
	assert.Equal(t, 1, roomAvailability.Quota.Value())

	id, err = vo.NewRoomAvailabilityID("reddison", "lux", timeutils.DateToTime(2024, 9, 22))
	assert.NoError(t, err)
	roomAvailability, err = storage.FindByID(context.Background(), *id)
	assert.Error(t, err)
	assert.Nil(t, roomAvailability)
}

func TestImMemoryRoomAvailabilityStorage_Save(t *testing.T) {
	storage := repository.NewImMemoryRoomAvailabilityStorage()

	quota, err := vo.NewQuotaFromInt(2)
	assert.NoError(t, err)

	roomAvailability, err := entity.NewRoomAvailability(
		vo.NewHotelIDFromString("foo"),
		vo.NewRoomIDFromString("bar"),
		timeutils.DateToTime(2024, 8, 20),
		*quota,
	)
	assert.NoError(t, err)
	err = storage.Save(context.Background(), roomAvailability)
	assert.NoError(t, err)

	id, err := vo.NewRoomAvailabilityID("foo", "bar", timeutils.DateToTime(2024, 8, 20))
	assert.NoError(t, err)
	foundRoomAvailability, err := storage.FindByID(context.Background(), *id)
	assert.NoError(t, err)
	assert.NotNil(t, foundRoomAvailability)
	assert.Equal(t, roomAvailability, foundRoomAvailability)

	err = storage.Save(context.Background(), roomAvailability)
	assert.Error(t, err)
	assert.Equal(t, exception.ErrRoomAvailabilityAlreadyExists, err)
}

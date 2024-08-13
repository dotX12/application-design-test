package repository_test

import (
	"applicationDesignTest/internal/adapter/storage/repository"
	"applicationDesignTest/internal/domain/entity"
	"applicationDesignTest/internal/domain/vo"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestImMemoryOrderStorage_SaveAndFindByID(t *testing.T) {
	storage := repository.NewImMemoryOrderStorage()

	orderID := vo.NewOrderID()

	email, err := vo.NewEmailFromString("user@example.com")
	assert.NoError(t, err)

	order := &entity.Order{
		ID:        orderID,
		HotelID:   vo.NewHotelIDFromString("reddison"),
		RoomID:    vo.NewRoomIDFromString("lux"),
		UserEmail: email,
		From:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		To:        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	err = storage.Create(context.Background(), order)
	assert.NoError(t, err)

	foundOrder, err := storage.FindByID(context.Background(), orderID)
	assert.NoError(t, err)
	assert.NotNil(t, foundOrder)
	assert.Equal(t, order, foundOrder)

	nonExistentOrder, err := storage.FindByID(context.Background(), vo.NewOrderID())
	assert.Error(t, err)
	assert.Nil(t, nonExistentOrder)

}

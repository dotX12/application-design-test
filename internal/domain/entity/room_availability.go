package entity

import (
	"applicationDesignTest/internal/domain/common/entity"
	"applicationDesignTest/internal/domain/exception"
	"applicationDesignTest/internal/domain/vo"
	"time"
)

type RoomAvailability struct {
	entity.AggregateRoot
	ID      vo.RoomAvailabilityID
	HotelID vo.HotelID
	RoomID  vo.RoomID
	Date    time.Time
	Quota   vo.Quota
}

func (ra *RoomAvailability) IsAvailable() bool {
	return ra.Quota.Value() > 0
}

func (ra *RoomAvailability) Reserve() error {
	if ra.IsAvailable() {
		quota, err := vo.NewQuotaFromInt(ra.Quota.Value() - 1)
		if err != nil {
			return err
		}
		ra.Quota = *quota
	}
	return nil
}

func NewRoomAvailability(hotelID vo.HotelID, roomID vo.RoomID, date time.Time, quota vo.Quota) (*RoomAvailability, error) {
	if quota.Value() < 0 {
		return nil, exception.ErrQuotaIsNegative
	}
	id, err := vo.NewRoomAvailabilityID(hotelID.String(), roomID.String(), date)
	if err != nil {
		return nil, err
	}
	return &RoomAvailability{
		ID:      *id,
		HotelID: hotelID,
		RoomID:  roomID,
		Date:    date,
		Quota:   quota,
	}, nil
}

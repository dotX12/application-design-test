package vo

import (
	"time"
)

type RoomAvailabilityID struct {
	value string
}

func (r RoomAvailabilityID) String() string {
	return r.value
}

func NewRoomAvailabilityID(
	hotelID string,
	roomID string,
	date time.Time,
) (*RoomAvailabilityID, error) {

	v := hotelID + "_" + roomID + "_" + date.Format("2006-01-02")
	return &RoomAvailabilityID{value: v}, nil
}

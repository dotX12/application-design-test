package request

import "time"

type CreateRoomAvailabilityRequest struct {
	RoomID  string    `json:"room_id"`
	HotelID string    `json:"hotel_id"`
	Date    time.Time `json:"date" format:"date-time"`
	Quota   int       `json:"quota"`
} // @name CreateRoomAvailabilityRequest

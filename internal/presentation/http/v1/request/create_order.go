package request

import "time"

type CreateOrderRequest struct {
	RoomID    string    `json:"room_id" example:"lux"`
	HotelID   string    `json:"hotel_id" example:"reddison"`
	UserEmail string    `json:"email" example:"foo@bar.com"`
	From      time.Time `json:"from" format:"date-time" example:"2024-09-12T00:00:00Z"`
	To        time.Time `json:"to" format:"date-time" example:"2024-09-15T00:00:00Z"`
} // @name CreateOrderRequest

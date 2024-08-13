package response

import (
	"applicationDesignTest/internal/domain/entity"
)

type OrderResponse struct {
	ID        string `json:"id"`
	RoomID    string `json:"room_id" example:"premium"`
	HotelID   string `json:"hotel_id" example:"resort"`
	UserEmail string `json:"email" example:"foo@bar.com"`
	From      Date   `json:"from" format:"date" swaggertype:"string" example:"2021-01-01"`
	To        Date   `json:"to" format:"date" swaggertype:"string" example:"2021-01-02"`
} // @name OrderResponse

func NewOrderResponseFromEntity(order *entity.Order) *OrderResponse {
	return &OrderResponse{
		ID:        order.ID.String(),
		RoomID:    order.RoomID.String(),
		HotelID:   order.HotelID.String(),
		UserEmail: order.UserEmail.String(),
		From:      Date{order.From},
		To:        Date{order.To},
	}
}

package response

type RoomAvailabilityCreatedResponse struct {
	ID string `json:"id"`
}

func NewRoomAvailabilityCreatedResponse(id string) *RoomAvailabilityCreatedResponse {
	return &RoomAvailabilityCreatedResponse{
		ID: id,
	}
}

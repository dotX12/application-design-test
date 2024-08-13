package vo

type RoomID struct {
	value string
}

func (r RoomID) String() string {
	return r.value
}

func NewRoomIDFromString(value string) RoomID {
	return RoomID{value: value}
}

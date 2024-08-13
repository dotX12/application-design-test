package vo

type HotelID struct {
	value string
}

func (h HotelID) String() string {
	return h.value
}

func NewHotelIDFromString(value string) HotelID {
	return HotelID{value: value}
}

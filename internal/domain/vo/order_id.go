package vo

import "github.com/google/uuid"

type OrderID struct {
	value string
}

func (o OrderID) String() string {
	return o.value
}

func NewOrderID() OrderID {
	return OrderID{value: uuid.New().String()}
}

func NewOrderIDFromString(value string) OrderID {
	return OrderID{value: value}
}

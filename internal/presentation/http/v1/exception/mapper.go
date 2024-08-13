package exception

import (
	"applicationDesignTest/internal/application/common/exception"
	"applicationDesignTest/pkg/timeutils"
	"errors"

	"github.com/valyala/fasthttp"
)

type ErrorCode struct {
	HTTPCode int
	Status   string
	Slug     string
	Message  string
}

func NewErrorCode(httpCode int, status, slug, message string) ErrorCode {
	return ErrorCode{
		HTTPCode: httpCode,
		Status:   status,
		Slug:     slug,
		Message:  message,
	}
}

var ErrCodeMap = map[error]ErrorCode{
	exception.ErrOrderNotFound: NewErrorCode(
		404,
		fasthttp.StatusMessage(fasthttp.StatusNotFound),
		"404_order_not_found",
		"The requested order was not found.",
	),
	exception.ErrOrderAlreadyExists: NewErrorCode(
		400,
		fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		"400_order_already_exists",
		"The requested order already exists.",
	),
	exception.ErrRoomAvailabilityNotFound: NewErrorCode(
		404,
		fasthttp.StatusMessage(fasthttp.StatusNotFound),
		"404_room_availability_not_found",
		"The requested room availability was not found.",
	),
	exception.ErrRoomAvailabilityAlreadyExists: NewErrorCode(
		400,
		fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		"400_room_availability_already_exists",
		"The requested room availability already exists.",
	),
	exception.ErrHotelRoomNotAvailable: NewErrorCode(
		400,
		fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		"400_hotel_room_not_available",
		"The hotel room is not available for selected dates.",
	),
	exception.ErrFindRoomAvailability: NewErrorCode(
		400,
		fasthttp.StatusMessage(fasthttp.StatusNotFound),
		"400_find_room_availability",
		"Failed to find room availability.",
	),
	exception.ErrSaveRoomAvailability: NewErrorCode(
		400,
		fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		"400_save_room_availability",
		"Failed to save room availability.",
	),
	exception.ErrSaveOrder: NewErrorCode(
		400,
		fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		"400_save_order",
		"Failed to save order.",
	),
	timeutils.ErrFromAfterTo: NewErrorCode(
		400,
		fasthttp.StatusMessage(fasthttp.StatusBadRequest),
		"400_from_after_to",
		"From date is after to date.",
	),
}

func GetErrorCodeFromError(err error) (*ErrorCode, bool) {
	for key, code := range ErrCodeMap {
		if errors.Is(err, key) {
			return &code, true
		}
	}

	return nil, false
}

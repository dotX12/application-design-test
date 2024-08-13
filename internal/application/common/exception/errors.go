package exception

import (
	domainErr "applicationDesignTest/internal/domain/exception"
	"fmt"
)

var ErrOrderNotFound = fmt.Errorf("%w: order not found", domainErr.ErrApplication)
var ErrOrderAlreadyExists = fmt.Errorf("%w: order already exists", domainErr.ErrApplication)

var ErrRoomAvailabilityNotFound = fmt.Errorf("%w: room availability not found", domainErr.ErrApplication)
var ErrRoomAvailabilityAlreadyExists = fmt.Errorf("%w: room availability already exists", domainErr.ErrApplication)

var ErrHotelRoomNotAvailable = fmt.Errorf("%w: hotel room is not available for selected dates", domainErr.ErrApplication)
var ErrFindRoomAvailability = fmt.Errorf("%w: failed to find room availability", domainErr.ErrApplication)

var ErrSaveRoomAvailability = fmt.Errorf("%w: failed to save room availability", domainErr.ErrApplication)

var ErrSaveOrder = fmt.Errorf("%w: failed to save order", domainErr.ErrApplication)

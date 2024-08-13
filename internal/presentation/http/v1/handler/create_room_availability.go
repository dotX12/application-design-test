package handler

import (
	"applicationDesignTest/internal/application/command"
	"applicationDesignTest/internal/domain/vo"
	"applicationDesignTest/internal/presentation/http/v1/exception"
	"applicationDesignTest/internal/presentation/http/v1/request"
	"applicationDesignTest/internal/presentation/http/v1/response"
	"applicationDesignTest/pkg/mediator"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// CreateRoomAvailabilityHandler
//
//	@Summary	Create new room availability
//	@Tags		rooms
//	@Accept		json
//	@Produce	json
//	@Param		request	body	CreateRoomAvailabilityRequest		true	"-"
//	@Success	200				{object}			OrderResponse
//	@Failure	400				{object}			ErrorResponse
//	@Failure	422				{object}			ErrorResponse
//	@Router		/v1/rooms [post]
func (h Handler) CreateRoomAvailabilityHandler(c *fiber.Ctx) error {
	h.observer.Logger.Trace().Ctx(c.UserContext()).Msg("Received request to CreateRoomAvailabilityHandler")

	var req request.CreateRoomAvailabilityRequest

	err := c.BodyParser(&req)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("CreateRoomAvailabilityHandler: error parsing request")
		return exception.HTTPRespondWithError(
			c,
			err,
			"422_malformed_request",
			"Unprocessable Entity",
			"The request cannot be processed due to malformed syntax",
			http.StatusUnprocessableEntity,
		)
	}

	cmd := &command.AddRoomAvailabilityCommand{
		HotelID: req.HotelID,
		RoomID:  req.RoomID,
		Date:    req.Date,
		Quota:   req.Quota,
	}
	send, err := mediator.Send[*command.AddRoomAvailabilityCommand, *vo.RoomAvailabilityID](c.UserContext(), cmd)
	if err != nil {
		if code, ok := exception.GetErrorCodeFromError(err); ok {
			return exception.HTTPRespondWithError(
				c,
				err,
				code.Slug,
				code.Status,
				code.Message,
				code.HTTPCode,
			)
		}
		return exception.HTTPRespondWithError(
			c,
			err,
			"400_unexpected_error",
			"Bad Request",
			"Unexpected Error",
			400,
		)
	}

	return c.Status(http.StatusCreated).JSON(response.NewRoomAvailabilityCreatedResponse(send.String()))
}

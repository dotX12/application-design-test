package handler

import (
	"applicationDesignTest/internal/application/command"
	"applicationDesignTest/internal/domain/entity"
	"applicationDesignTest/internal/presentation/http/v1/exception"
	"applicationDesignTest/internal/presentation/http/v1/request"
	"applicationDesignTest/internal/presentation/http/v1/response"
	"applicationDesignTest/pkg/mediator"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// CreateOrderHandler
//
//	@Summary	Create new order
//	@Tags		orders
//	@Accept		json
//	@Produce	json
//	@Param		request	body	CreateOrderRequest		true	"-"
//	@Success	200				{object}			OrderResponse
//	@Failure	400				{object}			ErrorResponse
//	@Failure	422				{object}			ErrorResponse
//	@Router		/v1/orders [post]
func (h Handler) CreateOrderHandler(c *fiber.Ctx) error {
	h.observer.Logger.Trace().Ctx(c.UserContext()).Msg("Received request to CreateOrderHandler")

	var req request.CreateOrderRequest

	err := c.BodyParser(&req)
	if err != nil {
		h.observer.Logger.Error().Err(err).Msg("CreateOrderHandler: error parsing request")
		return exception.HTTPRespondWithError(
			c,
			err,
			"422_malformed_request",
			"Unprocessable Entity",
			"The request cannot be processed due to malformed syntax",
			http.StatusUnprocessableEntity,
		)
	}

	cmd := &command.AddOrderCommand{
		HotelID:   req.HotelID,
		RoomID:    req.RoomID,
		UserEmail: req.UserEmail,
		From:      req.From,
		To:        req.To,
	}
	send, err := mediator.Send[*command.AddOrderCommand, *entity.Order](c.UserContext(), cmd)
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

	return c.Status(http.StatusCreated).JSON(response.NewOrderResponseFromEntity(send))
}

package handler

import (
	"applicationDesignTest/internal/presentation/http/v1/response"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// HealthcheckHandler
//
//	@Summary	Healthcheck
//	@Tags		healthcheck
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}			HealthCheckResponse
//	@Failure	400				{object}			ErrorResponse
//	@Failure	422				{object}			ErrorResponse
//	@Router		/v1/healthcheck [get]
func (h Handler) HealthcheckHandler(c *fiber.Ctx) error {
	h.observer.Logger.Trace().Ctx(c.UserContext()).Msg("Received request to HealthcheckHandler")

	return c.Status(http.StatusOK).JSON(response.HealthCheckResponse{Status: "OK"})
}

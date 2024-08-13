package exception

import (
	"applicationDesignTest/internal/presentation/http/v1/response"
	"github.com/gofiber/fiber/v2"
)

func HTTPRespondWithError(
	ctx *fiber.Ctx,
	err error,
	slug string,
	status string,
	message string,
	code int,
) error {
	var details string
	if err != nil {
		details = err.Error()
	}

	resp := response.ErrorResponse{
		Error: response.ErrorDetailsResponse{
			Slug:    slug,
			Message: message,
			Status:  status,
			Details: []response.ErrorDetailComponentResponse{
				{
					Field: []string{},
					Error: details,
				},
			},
		},
	}
	return ctx.Status(code).JSON(resp)
}

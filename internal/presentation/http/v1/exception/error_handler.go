package exception

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if err != nil {
		span := trace.SpanFromContext(ctx.UserContext())
		span.RecordError(err, trace.WithStackTrace(true))
		span.SetStatus(codes.Error, err.Error())

		return HTTPRespondWithError(
			ctx,
			err,
			"400_malformed_request",
			"Bad Request",
			"The request cannot be processed due to malformed syntax",
			code,
		)
	}
	return nil
}

package opentelemetry

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/propagation"
)

var traceParentHeader = "traceparent"

func TraceParentMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		res := c.Next()
		header := http.Header{}
		propagation.TraceContext{}.Inject(c.UserContext(), propagation.HeaderCarrier(header))
		traceParent := header.Get(traceParentHeader)
		c.Append(traceParentHeader, traceParent)
		return res
	}
}

package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Logger struct {
	zerolog.Logger
}

func NewZeroLogger(level zerolog.Level, hooks []zerolog.Hook) *Logger {
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()
	for _, hook := range hooks {
		logger = logger.Hook(hook)
	}
	return &Logger{logger}
}

type OTELHook struct{}

func (l OTELHook) Run(
	e *zerolog.Event,
	level zerolog.Level,
	message string,
) {
	ctx := e.GetCtx()

	if level == zerolog.NoLevel {
		return
	}
	if !e.Enabled() {
		return
	}

	if ctx == nil {
		return
	}
	span := trace.SpanFromContext(ctx)

	if !span.IsRecording() {
		return
	}
	{
		sCtx := span.SpanContext()
		if sCtx.HasTraceID() {
			e.Str("traceId", sCtx.TraceID().String())
		}
		if sCtx.HasSpanID() {
			e.Str("spanId", sCtx.SpanID().String())
		}
	}

	{
		attrs := make([]attribute.KeyValue, 0)

		logSeverityKey := attribute.Key("log.severity")
		logMessageKey := attribute.Key("log.message")
		attrs = append(attrs, logSeverityKey.String(level.String()), logMessageKey.String(message))

		span.AddEvent("log", trace.WithAttributes(attrs...))
		if level >= zerolog.ErrorLevel {
			span.SetStatus(codes.Error, message)
		}
	}
}

var DefaultLogger = NewZeroLogger(zerolog.TraceLevel, []zerolog.Hook{
	OTELHook{},
})

package tracer

import (
	"applicationDesignTest/pkg/logging"
	"context"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"google.golang.org/grpc/encoding/gzip"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type OpenTelemetry struct {
	Trace trace.Tracer
	Tp    *sdktrace.TracerProvider
}

func New(config Config, logger *logging.Logger) func(context.Context) error {
	ctx := context.Background()
	secureOption := otlptracegrpc.WithInsecure()
	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(config.Endpoint),
			otlptracegrpc.WithCompressor(gzip.Name),
		),
	)
	if err != nil {
		panic(err)
	}

	resources, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			attribute.String("service.name", config.ServiceName),
			attribute.String("service.version", "1.0.0"),
			attribute.String("library.language", "go"),
		))
	if err != nil {
		panic(err)
	}
	bsp := sdktrace.NewBatchSpanProcessor(
		exporter,
		sdktrace.WithMaxQueueSize(1000),
		sdktrace.WithMaxExportBatchSize(1000),
	)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithIDGenerator(xray.NewIDGenerator()),
		sdktrace.WithResource(resources),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.Error().Err(err).Msg("error in tracer")
	}))
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return bsp.Shutdown
}

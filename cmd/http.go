package main

import (
	"applicationDesignTest/config"
	"applicationDesignTest/docs"
	"applicationDesignTest/internal/adapter/storage/repository"
	"applicationDesignTest/internal/application/command"
	eventHandler "applicationDesignTest/internal/application/event"
	"applicationDesignTest/internal/application/pipeline"
	"applicationDesignTest/internal/application/query"
	"applicationDesignTest/internal/domain/entity"
	"applicationDesignTest/internal/domain/event"
	"applicationDesignTest/internal/domain/vo"
	"applicationDesignTest/internal/observability"
	"applicationDesignTest/internal/presentation/http/v1/exception"
	"applicationDesignTest/internal/presentation/http/v1/handler"
	"applicationDesignTest/internal/presentation/http/v1/middleware/opentelemetry"
	"applicationDesignTest/pkg/docs/rapidoc"
	"applicationDesignTest/pkg/docs/swagger"
	"applicationDesignTest/pkg/logging"
	"applicationDesignTest/pkg/logging/fiberzerolog"
	"applicationDesignTest/pkg/mediator"
	"applicationDesignTest/pkg/tracer"
	"context"
	"fmt"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recovery "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"html/template"
	"os"
	"os/signal"
	"syscall"
)

func makeMediator(logger *logging.Logger, traceWrapper *tracer.Wrapper) {
	observer := observability.New(logger, traceWrapper)
	orderStorage := repository.NewImMemoryOrderStorage()
	roomAvailabilityStorage := repository.NewImMemoryRoomAvailabilityStorage()
	addOrderHandler := command.NewAddOrderCommandHandler(
		observer,
		orderStorage,
		roomAvailabilityStorage,
	)
	getOrderByIDHandler := query.NewGetOrderByIDHandler(
		observer,
		orderStorage,
	)

	orderCreatedEventHandler := eventHandler.NewNotificationOrderCreatedHandler(observer)

	roomAvailabilityHandler := command.NewRoomAvailabilityCommandHandler(
		observer,
		roomAvailabilityStorage,
	)
	loggerPipeline := pipeline.NewRequestLoggerBehaviour(observer)
	err := mediator.RegisterRequestPipelineBehaviors(loggerPipeline)

	if err != nil {
		panic(err)
	}

	err = mediator.RegisterRequestHandler[*command.AddOrderCommand, *entity.Order](addOrderHandler)
	if err != nil {
		panic(err)
	}
	err = mediator.RegisterRequestHandler[*query.GetOrderByIDQuery, *entity.Order](getOrderByIDHandler)
	if err != nil {
		panic(err)
	}

	err = mediator.RegisterRequestHandler[*command.AddRoomAvailabilityCommand, *vo.RoomAvailabilityID](roomAvailabilityHandler)
	if err != nil {
		panic(err)
	}
	err = mediator.RegisterNotificationHandler[event.OrderCreated](orderCreatedEventHandler)
	if err != nil {
		panic(err)
	}
}

func setupMiddleware(
	cfg *config.Factory,
	observer *observability.Observability,
	app *fiber.App,

) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.API.OriginCORS,
		AllowCredentials: cfg.API.AllowCredentials,
	}))

	app.Use(
		otelfiber.Middleware(
			otelfiber.WithServerName(cfg.OTLP.ServiceName),
		),
	)
	app.Use(opentelemetry.TraceParentMiddleware())

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: observer.GetLogger(),
		Fields: []string{
			fiberzerolog.FieldIP,
			fiberzerolog.FieldUserAgent,
			fiberzerolog.FieldLatency,
			fiberzerolog.FieldStatus,
			fiberzerolog.FieldMethod,
			fiberzerolog.FieldURL,
			fiberzerolog.FieldError,
			fiberzerolog.FieldBytesReceived,
			fiberzerolog.FieldBytesSent,
		},
	}))

	app.Use(recovery.New(recovery.Config{EnableStackTrace: true}))
	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))
}

func setupRoutes(
	config *config.Factory,
	observer *observability.Observability,
	app *fiber.App,
) {
	v1Handler := handler.New(observer)

	v1Router := app.Group("/v1")
	v1Router.Get("/healthcheck", v1Handler.HealthcheckHandler)

	v1Router.Post("/rooms", v1Handler.CreateRoomAvailabilityHandler)
	v1Router.Post("/orders", v1Handler.CreateOrderHandler)
	v1Router.Get("/orders/:order_id", v1Handler.GetOrderByIDHandler)

	v1Router.Get("/swagger/*", swagger.New(swagger.Config{
		Title:                  "Booking Documentation",
		InstanceName:           "v1",
		URL:                    "/v1/swagger/openapi.json",
		DeepLinking:            false,
		DisplayRequestDuration: true,
		Presets: []template.JS{
			"SwaggerUIBundle.presets.apis",
			"SwaggerUIStandalonePreset",
		},
		Layout: "BaseLayout",
	}))
	v1Router.Get("/rapidoc/*", rapidoc.New(rapidoc.Config{
		Title:      "Booking Documentation",
		HeaderText: "Booking Documentation",
		SpecURL:    "swagger/openapi.json",

		LogoURL: "https://raw.githubusercontent.com/rapi-doc/RapiDoc/ebda9d7b3ac0d1b35ee4210c4a493a01567f4c87/docs/images/logo-outline.svg",
		Theme:   rapidoc.Theme_Light,
	}))
	docs.SwaggerInfov1.Host = fmt.Sprintf("%s:%s", config.API.Host, config.API.Port)

}

// @title						Booking Documentation
// @version						1.0
// @description					This is a sample swagger for Booking 2gis microservice.
// @host						localhost:8000
// @BasePath					/
func main() {
	cfg := config.New()

	logger := logging.NewZeroLogger(
		zerolog.TraceLevel,
		[]zerolog.Hook{
			logging.OTELHook{},
		},
	)

	cleanupTracer := tracer.New(
		tracer.Config{
			ServiceName:  cfg.OTLP.ServiceName,
			Endpoint:     cfg.OTLP.Endpoint,
			InSecureMode: cfg.OTLP.InSecureMode,
		},
		logger,
	)
	defer func() {
		if err := cleanupTracer(context.Background()); err != nil {
			logger.Error().Err(err).Msg("Error cleaning up tracer")
		}
	}()
	traceWrapper := tracer.NewWrapper(
		otel.Tracer(cfg.OTLP.ServiceName, trace.WithInstrumentationVersion("v0.1.0")),
	)
	observer := observability.New(logger, traceWrapper)

	observer.Logger.Info().
		Ctx(context.Background()).
		Str("host", cfg.API.Host).
		Str("port", cfg.API.Port).
		Msg("Starting HTTP server")

	makeMediator(logger, traceWrapper)

	app := fiber.New(fiber.Config{
		AppName:               "Booking microservice",
		ErrorHandler:          exception.ErrorHandler,
		DisableStartupMessage: true,
	})
	setupMiddleware(cfg, observer, app)
	setupRoutes(cfg, observer, app)

	routes := app.GetRoutes()
	for _, route := range routes {
		observer.Logger.Trace().
			Str("method", route.Method).
			Str("path", route.Path).
			Msg("Route added")
	}

	go func() {
		err := app.Listen(fmt.Sprintf("%s:%s", cfg.API.Host, cfg.API.Port))
		if err != nil {
			observer.Logger.Error().Err(err).Msgf("Error starting HTTP server")
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	sig := <-c

	observer.Logger.Info().Msgf("Received signal: %s", sig)
	observer.Logger.Info().Msg("Shutting down HTTP server...")
	_ = app.Shutdown()
	observer.Logger.Info().Msg("Running cleanup tasks...")
	mediator.ClearRequestRegistrations()
	mediator.ClearNotificationRegistrations()
	observer.Logger.Info().Msg("Cleanup tasks completed")

}

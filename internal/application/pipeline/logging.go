package pipeline

import (
	"applicationDesignTest/internal/observability"
	"applicationDesignTest/pkg/mediator"
	"context"
	"fmt"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

type RequestLoggerBehaviour struct {
	observer *observability.Observability
}

func NewRequestLoggerBehaviour(observer *observability.Observability) *RequestLoggerBehaviour {
	return &RequestLoggerBehaviour{observer: observer}
}

func (r *RequestLoggerBehaviour) Handle(ctx context.Context, request interface{}, next mediator.RequestHandlerFunc) (interface{}, error) {
	ctx, span := r.observer.Tracer.StartSpan(ctx, fmt.Sprintf("Execute: %s", GetTypeName(request)))
	defer span.End()

	r.observer.Logger.Trace().
		Ctx(ctx).
		Interface("request", request).
		Str("type", GetTypeName(request)).
		Msg("Received request for processing in mediator")

	response, err := next(ctx)
	if err != nil {
		r.observer.Logger.
			Error().
			Err(err).
			Ctx(ctx).
			Interface("request", request).
			Str("type", GetTypeName(request)).
			Msg("Error processing request in mediator")
		span.RecordError(err, trace.WithStackTrace(true))
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	r.observer.Logger.
		Trace().
		Ctx(ctx).
		Interface("request", request).
		Str("type", GetTypeName(request)).
		Interface("response", response).
		Msg("Request processed successfully in mediator")

	return response, nil
}

func GetTypeName(input interface{}) string {
	if input == nil {
		return ""
	}

	t := reflect.TypeOf(input)
	if t.Kind() != reflect.Ptr {
		return t.Name()
	}

	return t.Elem().Name()
}

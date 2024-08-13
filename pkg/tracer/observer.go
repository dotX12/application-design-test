package tracer

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	tr "go.opentelemetry.io/otel/trace"
)

type Wrapper struct {
	tracer tr.Tracer
}

func NewWrapper(tracer tr.Tracer) *Wrapper {
	return &Wrapper{
		tracer: tracer,
	}
}

func (tw *Wrapper) StartSpan(ctx context.Context, name string) (context.Context, tr.Span) {
	ctx, span := tw.tracer.Start(ctx, name)
	return ctx, span
}

func (tw *Wrapper) EndSpan(ctx context.Context, span tr.Span) {
	span.End()
}

func (tw *Wrapper) Error(ctx context.Context, err error, options ...tr.EventOption) {
	span := tr.SpanFromContext(ctx)
	span.RecordError(err, options...)
	span.SetStatus(codes.Error, err.Error())
}

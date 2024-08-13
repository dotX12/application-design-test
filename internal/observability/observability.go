package observability

import (
	"applicationDesignTest/pkg/logging"
	"applicationDesignTest/pkg/tracer"

	"github.com/rs/zerolog"
)

type Observability struct {
	Logger *logging.Logger
	Tracer *tracer.Wrapper
}

func (o Observability) GetLogger() *zerolog.Logger {
	return &o.Logger.Logger
}

func New(logger *logging.Logger, tracer *tracer.Wrapper) *Observability {
	return &Observability{
		Logger: logger,
		Tracer: tracer,
	}
}

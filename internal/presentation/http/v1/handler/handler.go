package handler

import "applicationDesignTest/internal/observability"

type Handler struct {
	observer *observability.Observability
}

func New(
	observer *observability.Observability,
) *Handler {
	return &Handler{observer: observer}
}

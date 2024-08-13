package query

import (
	"applicationDesignTest/internal/domain/entity"
	"applicationDesignTest/internal/domain/vo"
	"applicationDesignTest/internal/observability"
	"applicationDesignTest/internal/port/storage"
	"context"
)

type GetOrderByIDQuery struct {
	OrderID string
}

type GetOrderByIDHandler struct {
	observer     *observability.Observability
	orderStorage storage.OrderStorage
}

func NewGetOrderByIDHandler(
	observer *observability.Observability,
	orderStorage storage.OrderStorage,
) *GetOrderByIDHandler {
	return &GetOrderByIDHandler{
		observer:     observer,
		orderStorage: orderStorage,
	}
}

func (h GetOrderByIDHandler) Handle(ctx context.Context, query *GetOrderByIDQuery) (*entity.Order, error) {
	h.observer.Logger.Trace().Ctx(ctx).Msg("GetOrderByIDHandler.Handle")
	orderIDVO := vo.NewOrderIDFromString(query.OrderID)

	order, err := h.orderStorage.FindByID(ctx, orderIDVO)
	if err != nil {
		return nil, err
	}

	return order, nil
}

package storage

import (
	"applicationDesignTest/internal/domain/entity"
	"applicationDesignTest/internal/domain/vo"
	"context"
)

type OrderStorage interface {
	FindByID(ctx context.Context, orderID vo.OrderID) (*entity.Order, error)
	Create(ctx context.Context, order *entity.Order) error
}

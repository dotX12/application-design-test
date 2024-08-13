package repository

import (
	"applicationDesignTest/internal/application/common/exception"
	"applicationDesignTest/internal/domain/entity"
	"applicationDesignTest/internal/domain/vo"
	"context"
	"sync"
)

type ImMemoryOrderStorage struct {
	orders map[vo.OrderID]*entity.Order
	mu     sync.RWMutex
}

func (s *ImMemoryOrderStorage) Create(ctx context.Context, order *entity.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.orders[order.ID]; exists {
		return exception.ErrOrderAlreadyExists
	}

	s.orders[order.ID] = order

	return nil
}

func (s *ImMemoryOrderStorage) FindByID(ctx context.Context, orderID vo.OrderID) (*entity.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderID]
	if !ok {
		return nil, exception.ErrOrderNotFound
	}
	return order, nil
}

func NewImMemoryOrderStorage() *ImMemoryOrderStorage {
	return &ImMemoryOrderStorage{
		orders: map[vo.OrderID]*entity.Order{},
	}
}

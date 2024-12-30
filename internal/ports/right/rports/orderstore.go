package rports

import (
	"ecom-api/internal/application/core/types/entity"
)

type OrderStore interface {
	CreateOrder(entity.Order) (int, error)
	CreateOrderItem(entity.OrderItem) error
}

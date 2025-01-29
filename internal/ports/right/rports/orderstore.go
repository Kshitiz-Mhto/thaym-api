package rports

import (
	"ecom-api/internal/application/core/types/entity"
)

type OrderStore interface {
	CreateOrder(order entity.Order) (string, error)           // Create a new order and return its ID
	GetOrderByID(orderID string) (*entity.Order, error)       // Retrieve an order by its ID
	GetOrdersByUserID(userID string) ([]*entity.Order, error) // Retrieve all orders for a specific user
	UpdateOrder(order entity.Order) error                     // Update an existing order
	DeleteOrder(orderID string) error                         // Delete an order and its associated items

	CreateOrderItem(orderItem entity.OrderItem) error                   // Add an item to an order
	GetOrderItemsByOrderId(orderID string) ([]*entity.OrderItem, error) // Retrieve all items for a specific order
	DeleteOrderItem(orderItemID string) error                           // Delete a specific order item

}

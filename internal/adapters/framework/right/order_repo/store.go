package order

import (
	"database/sql"
	"ecom-api/internal/application/core/types/entity"
	"fmt"
	"log"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (store *Store) CreateOrder(order entity.Order) (string, error) {
	_, err := store.db.Exec("INSERT INTO orders (userId, total, subtotal, status, paymentStatus, paymentMethod, address, currency) VALUES (?,?,?,?,?,?,?,?)", order.UserID, order.Total, order.Subtotal, order.Status, order.PaymentStatus, order.PaymentMethod, order.Address, order.Currency)
	if err != nil {
		return "", err
	}

	var uuid string
	err = store.db.QueryRow("SELECT id FROM orders ORDER BY createdAt DESC LIMIT 1").Scan(&uuid)
	if err != nil {
		return "", fmt.Errorf("failed to fetch order UUID: %w", err)
	}

	return uuid, nil
}

func (store *Store) GetOrderByID(orderID string) (*entity.Order, error) {
	rows, err := store.db.Query("SELECT * FROM orders WHERE id = ?", orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order by id: %w", err)
	}
	defer rows.Close()

	order := new(entity.Order)

	for rows.Next() {
		order, err = ScanRowsIntoOrder(rows)
		if err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return order, nil
}

func (store *Store) GetOrdersByUserID(userID string) ([]*entity.Order, error) {
	rows, err := store.db.Query("SELECT * FROM orders WHERE userId = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}

	defer rows.Close()

	var orders []*entity.Order

	for rows.Next() {
		order, err := ScanRowsIntoOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (store *Store) UpdateOrder(order entity.Order) error {
	_, err := store.db.Exec("UPDATE orders SET total = ?, subtotal = ?, status = ?, paymentStatus = ?, paymentMethod = ?, address = ?, currency = ?, updatedAt = NOW() WHERE id = ?", order.Total, order.Subtotal, order.Status, order.PaymentStatus, order.PaymentMethod, order.Address, order.Currency, order.ID)
	if err != nil {
		return err
	}
	return err
}

func (store *Store) DeleteOrder(orderID string) error {
	_, err := store.db.Exec("DELETE FROM orders WHERE id=?", orderID)
	if err != nil {
		return err
	}
	return err
}

func (store *Store) CreateOrderItem(orderitem entity.OrderItem) error {

	_, err := store.db.Exec("INSERT INTO orderitems (orderId, productId, productName, quantity, price, totalPrice, subTotal, currency, discount, tax) VALUES (?,?,?,?,?,?,?,?,?,?)", orderitem.OrderID, orderitem.ProductID, orderitem.ProductName, orderitem.Quantity, orderitem.Price, orderitem.TotalPrice, orderitem.Subtotal, orderitem.Currency, orderitem.Discount, orderitem.Tax)

	if err != nil {
		log.Println(err)
		return err
	}

	return err
}

func (store *Store) GetOrderItemsByOrderId(orderId string) ([]*entity.OrderItem, error) {
	rows, err := store.db.Query("SELECT * FROM orderitems WHERE orderId = ?", orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to query orderitems: %w", err)
	}

	defer rows.Close()

	var orderitems []*entity.OrderItem

	for rows.Next() {
		orderitem, err := ScanRowsIntoOrderItem(rows)
		if err != nil {
			return nil, err
		}
		orderitems = append(orderitems, orderitem)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orderitems, nil
}

func (store *Store) DeleteOrderItem(orderItemId string) error {
	_, err := store.db.Exec("DELETE FROM orderitems WHERE id=?", orderItemId)
	if err != nil {
		return err
	}
	return err
}

func ScanRowsIntoOrder(rows *sql.Rows) (*entity.Order, error) {
	order := new(entity.Order)

	err := rows.Scan(
		&order.ID,
		&order.UserID,
		&order.Total,
		&order.Subtotal,
		&order.Status,
		&order.PaymentStatus,
		&order.PaymentMethod,
		&order.Address,
		&order.Currency,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func ScanRowsIntoOrderItem(rows *sql.Rows) (*entity.OrderItem, error) {
	orderItem := new(entity.OrderItem)

	err := rows.Scan(
		&orderItem.ID,
		&orderItem.OrderID,
		&orderItem.ProductID,
		&orderItem.ProductName,
		&orderItem.Quantity,
		&orderItem.Price,
		&orderItem.TotalPrice,
		&orderItem.Subtotal,
		&orderItem.Currency,
		&orderItem.Discount,
		&orderItem.Tax,
		&orderItem.CreatedAt,
		&orderItem.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return orderItem, nil
}

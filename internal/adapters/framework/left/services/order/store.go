package order

import (
	"database/sql"
	"ecom-api/internal/application/core/types/entity"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (store *Store) CreateOrder(order entity.Order) (int, error) {
	_, err := store.db.Exec("INSERT INTO orders (userId, total, sub, tax, discout, status, paymentStatus, paymentMethod, address, currency) VALUES (?,?,?,?,?,?,?,?,?,?)", order.UserID, order.Total, order.Subtotal, order.Tax, order.Discount, order.Status, order.PaymentStatus, order.PaymentMethod, order.Address, order.Currency)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (store *Store) CreateOrderItem(orderitem entity.OrderItem) error {
	_, err := store.db.Exec("INSERT INTO orderitems (orderId, productId, productName, quatity, price, totalPrice, currency, discount, tax) VALUES (?,?,?,?,?,?,?,?,?)", orderitem.OrderID, orderitem.ProductID, orderitem.ProductName, orderitem.Quantity, orderitem.Price, orderitem.TotalPrice, orderitem.Currency, orderitem.Discount, orderitem.Tax)

	if err != nil {
		return err
	}

	return err
}

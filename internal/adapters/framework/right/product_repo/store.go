package product_repo

import (
	"database/sql"
	"ecom-api/internal/application/core/types/entity"
	"ecom-api/internal/application/core/types/entity/payloads"
	"fmt"
	"strings"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateProduct(product payloads.CreateProductPayload) error {
	_, err := s.db.Exec("INSERT INTO products(productId, name, description, image, price, currency, quantity, category, tags, isActive) VALUES (?,?,?,?,?,?,?,?,?,?)", product.ProductId, product.Name, product.Description, product.Image, product.Price, product.Currency, product.Quantity, product.Category, product.Tags, product.IsActive)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetProductByID(id string) (*entity.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERE productId =?", id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	product := new(entity.Product)
	for rows.Next() {
		product, err = scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Store) GetProductsByID(ids []string) ([]entity.Product, error) {

	if len(ids) == 0 {
		return nil, fmt.Errorf("slice of ids is empty")
	}

	placeholder := strings.Repeat(",?", len(ids)-1)
	// If productIDs contains [1, 2, 3], the query becomes:
	query := fmt.Sprintf("SELECT * FROM products WHERE productId IN (?%s)", placeholder)

	args := make([]interface{}, len(ids))
	for i, v := range args {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Initialize the products slice
	products := []entity.Product{}

	for rows.Next() {
		product, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Store) GetAllProducts() ([]*entity.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := make([]*entity.Product, 0)

	for rows.Next() {
		product, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Store) UpdateProduct(product entity.Product) error {
	_, err := s.db.Exec("UPDATE products SET name = ?, description = ?, image = ?, price = ?, currency = ?, quantity = ?, category = ?, tags = ?, isActive = ?, updatedAt = CURRENT_TIMESTAMP WHERE productId = ?", product.Name, product.Description, product.Image, product.Price, product.Currency, product.Quantity, product.Category, product.Tags, product.IsActive, product.ProductId)

	if err != nil {
		return err
	}
	return nil
}

func scanRowsIntoProduct(rows *sql.Rows) (*entity.Product, error) {
	product := new(entity.Product)
	err := rows.Scan(
		&product.ProductId,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Currency,
		&product.Quantity,
		&product.Category,
		&product.Tags,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// // If the Tags field is a valid JSON array in the database, unmarshal it
	// if product.Tags != nil {
	// 	if err = json.Unmarshal(product.Tags, &product.Tags); err != nil {
	// 		return nil, fmt.Errorf("failed to unmarshal tags")
	// 	}
	// }

	return product, err
}

// ActivateProduct implements rports.ProductStore.
func (s *Store) ActivateProduct(id string) error {
	panic("unimplemented")
}

// DeactivateProduct implements rports.ProductStore.
func (s *Store) DeactivateProduct(id string) error {
	panic("unimplemented")
}

// DecreaseProductStock implements rports.ProductStore.
func (s *Store) DecreaseProductStock(id string, quantity int) error {
	panic("unimplemented")
}

// DeleteProductByID implements rports.ProductStore.
func (s *Store) DeleteProductByID(id string) error {
	panic("unimplemented")
}

// GetProductsByCategory implements rports.ProductStore.
func (s *Store) GetProductsByCategory(category string) ([]entity.Product, error) {
	panic("unimplemented")
}

// IncreaseProductStock implements rports.ProductStore.
func (s *Store) IncreaseProductStock(id string, quantity int) error {
	panic("unimplemented")
}

// SearchProducts implements rports.ProductStore.
func (s *Store) SearchProducts(query string) ([]entity.Product, error) {
	panic("unimplemented")
}

// UpdateProductQuantity implements rports.ProductStore.
func (s *Store) UpdateProductQuantity(id string, quantity int) error {
	panic("unimplemented")
}

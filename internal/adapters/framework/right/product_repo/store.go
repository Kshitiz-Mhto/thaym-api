package product_repo

import (
	"database/sql"
	"ecom-api/internal/application/core/types/entity"
	"ecom-api/internal/application/core/types/entity/payloads"
	"encoding/json"
	"errors"
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
	tagsJSON, errr := json.Marshal(product.Tags)
	if errr != nil {
		return errr
	}
	_, err := s.db.Exec("INSERT INTO products(name, description, image, price, currency, quantity, category, tags, isActive) VALUES (?,?,?,?,?,?,?,?,?)", product.Name, product.Description, product.Image, product.Price, product.Currency, product.Quantity, product.Category, tagsJSON, product.IsActive)
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

func (s *Store) GetProductsByIDs(ids []string) ([]entity.Product, error) {

	if len(ids) == 0 {
		return nil, fmt.Errorf("slice of ids is empty")
	}

	// Build the query with the proper number of placeholders for IN clause
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf("SELECT * FROM products WHERE productId IN (%s)", strings.Join(placeholders, ","))

	args := make([]interface{}, len(ids))
	for i, v := range ids {
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

	return products, nil
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

func (s *Store) DeleteProductByID(productId string) error {
	_, err := s.db.Exec("DELETE FROM products WHERE productId = ?", productId)

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateProductByID(productId string) error {
	var product payloads.CreateProductPayload

	_, err := s.db.Exec("UPDATE products SET name = ?, description = ?, image = ?, price = ?, currency = ?, quantity = ?, category = ?, tags = ?, isActive = ?, updatedAt = CURRENT_TIMESTAMP WHERE productId = ?", product.Name, product.Description, product.Image, product.Price, product.Currency, product.Quantity, product.Category, product.Tags, product.IsActive, productId)

	if err != nil {
		return err
	}
	return nil
}

// ActivateProduct implements rports.ProductStore.
func (s *Store) ActivateProduct(id string) error {
	result, err := s.db.Exec("UPDATE products SET isActive = ? WHERE productId = ?", true, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no product found with the given ID")
	}

	return nil
}

// DeactivateProduct implements rports.ProductStore.
func (s *Store) DeactivateProduct(id string) error {
	result, err := s.db.Exec("UPDATE products SET isActive = ? WHERE id = ?", false, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no product found with the given ID")
	}

	return nil
}

func (s *Store) DecreaseProductStock(id string, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	var currentQuantity int
	err := s.db.QueryRow("SELECT quantity FROM products WHERE productId = ?", id).Scan(&currentQuantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("product not found")
		}
		return err
	}

	if currentQuantity < quantity {
		return fmt.Errorf("not enough stock: available %d, requested %d", currentQuantity, quantity)
	}

	_, err = s.db.Exec("UPDATE products SET quantity = ? WHERE productId = ?", currentQuantity-quantity, id)
	if err != nil {
		return err
	}

	return nil
}

// GetProductsByCategory implements rports.ProductStore.
func (s *Store) GetProductsByCategory(category string) ([]*entity.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERE category =?", category)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []*entity.Product

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

// IncreaseProductStock implements rports.ProductStore.
func (s *Store) IncreaseProductStock(id string, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	var currentQuantity int
	err := s.db.QueryRow("SELECT quantity FROM products WHERE productId = ?", id).Scan(&currentQuantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("product not found")
		}
		return err
	}

	_, err = s.db.Exec("UPDATE products SET quantity = ? WHERE productId = ?", currentQuantity+quantity, id)
	if err != nil {
		return err
	}

	return nil
}

// SearchProducts implements rports.ProductStore.
func (s *Store) SearchProducts(query string) ([]*entity.Product, error) {
	if query = strings.TrimSpace(query); query == "" {
		return nil, errors.New("search query cannot be empty")
	}

	query = strings.ToLower(query)
	products, err := s.GetAllProducts()
	if err != nil {
		return nil, err
	}

	results := []*entity.Product{}

	for _, product := range products {
		if containsIgnoreCase(product.Name, query) ||
			containsIgnoreCase(product.Description, query) ||
			containsIgnoreCase(product.Category, query) ||
			tagsContain(product.Tags, query) {
			results = append(results, product)
		}
	}

	if len(results) == 0 {
		return nil, errors.New("no products found")
	}

	return results, nil
}

// UpdateProductQuantity implements rports.ProductStore.
func (s *Store) UpdateProductQuantity(id string, quantity int) error {
	result, err := s.db.Exec("UPDATE products SET quantity = ? WHERE id = ?", quantity, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no product found with the given ID")
	}

	return nil
}

func (s *Store) GetProductsByStoreOwner(storeId string) ([]*entity.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERE storeId =?", storeId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []*entity.Product
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

func containsIgnoreCase(source, query string) bool {
	return strings.Contains(strings.ToLower(source), query)
}

func tagsContain(tagsData []byte, query string) bool {
	var tags []string
	if err := json.Unmarshal(tagsData, &tags); err != nil {
		return false
	}
	for _, tag := range tags {
		if containsIgnoreCase(tag, query) {
			return true
		}
	}
	return false
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

	// If the Tags field is a valid JSON array in the database, unmarshal it
	// if len(product.Tags) > 0 {
	// 	if err = json.Unmarshal(product.Tags, &tags); err != nil {
	// 		return nil, fmt.Errorf("failed to unmarshal tags")
	// 	}
	// }

	return product, err
}

package rports

import (
	"ecom-api/internal/application/core/types/entity"
	"ecom-api/internal/application/core/types/entity/payloads"
)

type ProductStore interface {
	// Basic CRUD operations
	GetProductByID(id string) (*entity.Product, error)         // Fetch a single product by its unique ID
	GetProductsByIDs(ids []string) ([]entity.Product, error)   // Fetch multiple products by their IDs
	GetAllProducts() ([]*entity.Product, error)                // Fetch all products
	CreateProduct(payload payloads.CreateProductPayload) error // Create a new product
	UpdateProduct(product entity.Product) error                // Update product details
	DeleteProductByID(id string) error                         // Delete a product by its ID
	UpdateProductByID(id string) error                         // Update a product by its ID

	// Filtering and searching
	GetProductsByCategory(category string) ([]entity.Product, error) // Fetch products by category
	SearchProducts(query string) ([]entity.Product, error)           // Search products by name, description, or tags

	// Inventory management
	UpdateProductQuantity(id string, quantity int) error // Update the quantity of a product
	IncreaseProductStock(id string, quantity int) error  // Increase product stock
	DecreaseProductStock(id string, quantity int) error  // Decrease product stock

	// Product activation
	ActivateProduct(id string) error   // Mark a product as active
	DeactivateProduct(id string) error // Mark a product as inactive
}

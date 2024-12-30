package rports

import (
	"ecom-api/internal/application/core/types/entity"
	"ecom-api/internal/application/core/types/entity/payloads"
)

type ProductStore interface {
	GetProductByID(id int) (*entity.Product, error)
	GetProductsByID(ids []int) ([]entity.Product, error)
	GetProducts() ([]*entity.Product, error)
	CreateProduct(payloads.CreateProductPayload) error
	UpdateProduct(entity.Product) error
}

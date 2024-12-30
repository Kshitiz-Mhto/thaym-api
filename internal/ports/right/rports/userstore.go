package rports

import (
	"ecom-api/internal/application/core/types/entity"
)

type UserStore interface {
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id int) (*entity.User, error)
	CreateUser(entity.User) error
}

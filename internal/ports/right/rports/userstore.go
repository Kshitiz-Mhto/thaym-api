package rports

import (
	"ecom-api/internal/application/core/types/entity"
)

type UserStore interface {
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id string) (*entity.User, error)
	CreateUser(user entity.User) error
}

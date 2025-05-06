package user

import (
	"context"
)

type Storage interface {
	GetUserCurrencies(ctx context.Context) ([]Currency, error)
}

type User struct {
	storage Storage
}

func New(storage Storage) *User {
	return &User{
		storage: storage,
	}
}

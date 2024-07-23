package service

import (
	"context"
	"cruda-app/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredential(ctx context.Context, email, password string) (domain.User, error)
}

type Users struct {
	repo UserRepository
}

func NewUsers(repo UserRepository) *Users {
	return &Users{
		repo: repo,
	}
}

func SingUp(ctx context.Context, inp domain.SingUpInput) error {
	return nil
}

func SingIn(ctx context.Context, email, password string) (string, error) {
	return "", nil
}

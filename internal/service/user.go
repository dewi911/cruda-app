package service

import (
	"context"
	"cruda-app/internal/domain"
	"time"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredential(ctx context.Context, email, password string) (domain.User, error)
}

type Users struct {
	repo   UserRepository
	hasher PasswordHasher
}

func NewUsers(repo UserRepository, hasher PasswordHasher) *Users {
	return &Users{
		repo:   repo,
		hasher: hasher,
	}
}

func (u *Users) SingUp(ctx context.Context, inp domain.SingUpInput) error {
	password, err := u.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return u.repo.Create(ctx, user)
}

func (u *Users) SingIn(ctx context.Context, email, password string) (string, error) {
	return "", nil
}

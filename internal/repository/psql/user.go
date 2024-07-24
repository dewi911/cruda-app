package psql

import (
	"context"
	"cruda-app/internal/domain"
	"database/sql"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db}
}

func (u *Users) Create(ctx context.Context, user domain.User) error {
	_, err := u.db.Exec("INSERT INTO users (name, email, password, registered_at) VALUES ($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, user.RegisteredAt)

	return err
}

func (u *Users) GetByCredential(ctx context.Context, email, password string) (domain.User, error) {
	return domain.User{}, nil
}

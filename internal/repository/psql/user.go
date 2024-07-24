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

func (r *Users) Create(ctx context.Context, user domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (name, email, password, registered_at) VALUES ($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, user.RegisteredAt)

	return err
}

func (r *Users) GetByCredential(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, name, email, registered_at FROM users WHERE email=$1 AND password=$2", email, password).
		Scan(&user.ID, &user.Name, &user.Email, &user.RegisteredAt)
	
	return user, err
}

package psql

import (
	"context"
	"database/sql"
	"github.com/dewi911/cruda-app/internal/domain"
)

type Tokens struct {
	db *sql.DB
}

func NewTokens(db *sql.DB) *Tokens {
	return &Tokens{db: db}
}

func (r *Tokens) Create(ctx context.Context, token domain.RefreshSession) error {
	_, err := r.db.Exec("INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		token.UserID, token.Token, token.ExpiresAt)

	return err
}

func (r *Tokens) Get(ctx context.Context, token string) (domain.RefreshSession, error) {
	var t domain.RefreshSession
	err := r.db.QueryRow("SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token=$1", token).
		Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt)
	if err != nil {
		return t, err
	}

	_, err = r.db.Exec("DELETE FROM refresh_tokens WHERE user_id=$1", t.UserID)
	if err != nil {
		return t, err
	}

	return t, err
}

func (r *Tokens) Delete(ctx context.Context, id int, token domain.RefreshSession) error {
	_, err := r.db.Exec("DELETE FROM refresh_tokens WHERE id=$3",
		token.Token, token.ExpiresAt, id)

	return err
}

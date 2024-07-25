package domain

import "errors"

var (
	ErrBookNotFound        = errors.New("Book not found")
	ErrRefreshTokenExpired = errors.New("Refresh token expired")
)

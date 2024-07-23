package domain

import (
	"github.com/go-playground/validator/v10"
	"time"
)

var validate *validator.Validate

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

type SingUpInput struct {
	Name     string `json:"name" validate:"required,gte=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (i SingUpInput) Validate() error {
	return validate.Struct(i)
}

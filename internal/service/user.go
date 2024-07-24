package service

import (
	"context"
	"cruda-app/internal/domain"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"strconv"
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

	hmaSecret []byte
}

func NewUsers(repo UserRepository, hasher PasswordHasher, secret []byte) *Users {
	return &Users{
		repo:      repo,
		hasher:    hasher,
		hmaSecret: secret,
	}
}

func (s *Users) SingUp(ctx context.Context, inp domain.SingUpInput) error {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return s.repo.Create(ctx, user)
}

func (s *Users) SingIn(ctx context.Context, inp domain.SingInInput) (string, error) {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return "", err
	}

	user, err := s.repo.GetByCredential(ctx, inp.Email, password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
	})

	return token.SignedString(s.hmaSecret)
}

func (s *Users) ParseToken(ctx context.Context, tokenString string) (int64, error) {
	t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.hmaSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return int64(id), nil
}

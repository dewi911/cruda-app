package service

import (
	"context"
	"cruda-app/internal/domain"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"math/rand"
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

type TokensRepository interface {
	Create(ctx context.Context, user domain.RefreshToken) error
	Get(ctx context.Context, token string) (domain.RefreshToken, error)
}

type Users struct {
	repo      UserRepository
	tokenRepo TokensRepository
	hasher    PasswordHasher

	hmaSecret []byte
	tokenTtl  time.Duration
}

func NewUsers(repo UserRepository, tokenRepo TokensRepository, hasher PasswordHasher, secret []byte, ttl time.Duration) *Users {
	return &Users{
		repo:      repo,
		tokenRepo: tokenRepo,
		hasher:    hasher,
		hmaSecret: secret,
		tokenTtl:  ttl,
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

func (s *Users) SingIn(ctx context.Context, inp domain.SingInInput) (string, string, error) {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetByCredential(ctx, inp.Email, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", domain.ErrUserNotFound
		}
		return "", "", err
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(s.tokenTtl).Unix(),
	})

	accessToken, err := t.SignedString(s.hmaSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := s.tokenRepo.Create(ctx, domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
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

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

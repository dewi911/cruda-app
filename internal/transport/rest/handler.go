package rest

import (
	"context"
	"cruda-app/internal/domain"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Books interface {
	Create(ctx context.Context, book domain.Book) error
	GetByID(ctx context.Context, id int64) (domain.Book, error)
	GetAll(ctx context.Context) ([]domain.Book, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, inp domain.UpdateBookInput) error
}

type User interface {
	SingUp(ctx context.Context, user domain.SingUpInput) error
	SingIn(ctx context.Context, inp domain.SingInInput) (string, string, error)
	ParseToken(ctx context.Context, accessToken string) (int64, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

type Handler struct {
	booksService Books
	usersService User
}

func NewHandler(books Books, users User) *Handler {
	return &Handler{
		booksService: books,
		usersService: users,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/sing-up", h.SingUp).Methods(http.MethodPost)
		auth.HandleFunc("/sing-in", h.SingIn).Methods(http.MethodGet)
		auth.HandleFunc("/refresh", h.refresh).Methods(http.MethodGet)
	}

	books := r.PathPrefix("/books").Subrouter()
	{
		books.Use(h.authMiddleware)

		books.HandleFunc("", h.createBook).Methods(http.MethodPost)
		books.HandleFunc("", h.getAllBooks).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.getBookByID).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.deleteBook).Methods(http.MethodDelete)
		books.HandleFunc("/{id:[0-9]+}", h.updateBook).Methods(http.MethodPut)
	}

	return r
}

func getIdFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id can't be 0")
	}

	return id, nil
}

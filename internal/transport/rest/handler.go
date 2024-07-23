package rest

import (
	"context"
	"cruda-app/internal/domain"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
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

type Handler struct {
	booksService Books
}

func NewHandler(books Books) *Handler {
	return &Handler{
		booksService: books,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	books := r.PathPrefix("/books").Subrouter()
	{
		books.HandleFunc("", h.createBook).Methods(http.MethodPost)
		books.HandleFunc("", h.getAllBooks).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.getBookByID).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.deleteBook).Methods(http.MethodDelete)
		books.HandleFunc("/{id:[0-9]+}", h.updateBook).Methods(http.MethodPut)
	}

	return r
}

func (h *Handler) getBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getBookByID",
			"problem": "getting id from request",
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := h.booksService.GetByID(context.TODO(), id)
	if err != nil {
		if errors.Is(err, domain.ErrBookNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.WithFields(log.Fields{
			"handler": "getBookByID",
			"problem": "getting book by id",
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(book)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getBookByID",
			"problem": "marshalling book",
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"problem": "reading request body",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book domain.Book
	if err := json.Unmarshal(reqBytes, &book); err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"problem": "unmarshalling request body",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.booksService.Create(context.TODO(), book)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "createBook",
			"problem": "creating book",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.booksService.GetAll(context.TODO())
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getAllBooks",
			"problem": "getting all books",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(books)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getAllBooks",
			"problem": "marshalling books",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteBook",
			"problem": "getting id from request",
		}).Error(err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.booksService.Delete(context.TODO(), id)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteBook",
			"problem": "deleting book",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"problem": "getting id from request",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"problem": "reading request body",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp domain.UpdateBookInput
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"problem": "unmarshalling request body",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.booksService.Update(context.TODO(), id, inp)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateBook",
			"problem": "updating book",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

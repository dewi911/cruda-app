package rest

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type CtxValue int

const (
	ctxUserID CtxValue = iota
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"uri":    r.RequestURI,
		}).Info()
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)
		if err != nil {
			logError("authMiddleware", "token from request failed", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, err := h.usersService.ParseToken(r.Context(), token)
		if err != nil {
			logError("authMiddleware", "token parsing failed", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserID, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}

func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty Authorization header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}

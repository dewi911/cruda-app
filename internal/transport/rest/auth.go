package rest

import (
	"cruda-app/internal/domain"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (h *Handler) SingUp(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logError("SingUp", "reading request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp domain.SingUpInput
	if err := json.Unmarshal(reqBytes, &inp); err != nil {
		logError("SingUp", "unmarshalling request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := inp.Validate(); err != nil {
		logError("SingUp", "validation request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.usersService.SingUp(r.Context(), inp)
	if err != nil {
		logError("SingUp", "singing up", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SingIn(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logError("SingIn", "reading request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp domain.SingInInput
	if err := json.Unmarshal(reqBytes, &inp); err != nil {
		logError("SingIn", "unmarshalling request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := inp.Validate(); err != nil {
		logError("SingIn", "validation request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.usersService.SingIn(r.Context(), inp)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			handleNotFoundError(w, err)
			return
		}

		logError("SingIn", "token sing-in", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]string{
		"token": accessToken,
	})
	if err != nil {
		logError("SingIn", "marshalling response body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("refresh_token=%s; HttpOnly", refreshToken))
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		logError("refresh", "getting refresh token", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logrus.Info("%s", cookie.Value)

	accessToken, refreshToken, err := h.usersService.RefreshTokens(r.Context(), cookie.Value)
	if err != nil {
		logError("refresh", "getting refresh token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]string{
		"token": accessToken,
	})

	if err != nil {
		logError("refresh", "marshalling response body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("refresh_token=%s; HttpOnly", refreshToken))
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func handleNotFoundError(w http.ResponseWriter, err error) {
	response, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(response)
}

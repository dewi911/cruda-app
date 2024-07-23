package rest

import (
	"cruda-app/internal/domain"
	"encoding/json"
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

}

package verify

import (
	"advancedgo/internal/email/storage"
	"advancedgo/pkg/response"
	"errors"
	"net/http"
)

type Verifyhandler struct {
	storage *storage.Storage
}

func NewVerifyHandler(router *http.ServeMux, storage *storage.Storage) *Verifyhandler {
	h := &Verifyhandler{
		storage: storage,
	}
	router.HandleFunc("GET /verify/{hash}", h.Verify())
	return h
}

func (h *Verifyhandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")

		emailAddr, ok := h.storage.Pop(hash)
		if !ok {
			response.ErrorResponse(w, http.StatusNotFound, errors.New("invalid or expired hash"))
			return
		}

		response.OkResponse(w, http.StatusOK, VerifyResponse{Message: "verified", Email: emailAddr})
	}
}

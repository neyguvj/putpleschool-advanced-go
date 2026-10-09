package verify

import (
	"advancedgo/internal/email/storage"
	"advancedgo/pkg/response"
	"errors"
	"log"
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

		emailAddr, ok, err := h.storage.Pop(hash)
		if err != nil {
			log.Printf("storage pop: %v", err)
			response.ErrorResponse(w, http.StatusInternalServerError, errors.New("failed to verify hash"))
			return
		}
		if !ok {
			response.ErrorResponse(w, http.StatusNotFound, errors.New("invalid or expired hash"))
			return
		}

		response.OkResponse(w, http.StatusOK, VerifyResponse{Verified: true, Email: emailAddr})
	}
}

package send

import (
	"advancedgo/internal/email/config"
	"advancedgo/internal/email/storage"
	"advancedgo/pkg/request"
	"advancedgo/pkg/response"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net"
	"net/http"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type SendHandler struct {
	cfg     *config.Config
	storage *storage.Storage
}

func NewSendHandler(router *http.ServeMux, cfg *config.Config, storage *storage.Storage) *SendHandler {
	h := &SendHandler{
		cfg:     cfg,
		storage: storage,
	}
	router.HandleFunc("POST /send", h.Send())
	return h
}

func (h *SendHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := request.DecodeAndValidate[SendRequest](w, r)
		if err != nil {
			response.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		hash, err := newHash()
		if err != nil {
			response.ErrorResponse(w, http.StatusInternalServerError, errors.New("failed to generate hash"))
			return
		}

		if err := h.storage.Put(hash, req.Email); err != nil {
			log.Printf("storage put: %v", err)
			response.ErrorResponse(w, http.StatusInternalServerError, errors.New("failed to save hash"))
			return
		}

		host, _, err := net.SplitHostPort(h.cfg.SmtpAddr)
		if err != nil {
			response.ErrorResponse(w, http.StatusInternalServerError, errors.New("invalid SMTP_ADDR"))
			return
		}

		e := email.NewEmail()
		e.From = h.cfg.Email
		e.To = []string{req.Email}
		e.Subject = "Verify your email"
		e.Text = []byte("Please confirm your email: http://localhost:8080/verify/" + hash)

		auth := smtp.PlainAuth("", h.cfg.Email, h.cfg.Password, host)
		if err := e.Send(h.cfg.SmtpAddr, auth); err != nil {
			log.Printf("send email to %s: %v", req.Email, err)
			response.ErrorResponse(w, http.StatusBadGateway, errors.New("failed to send email"))
			return
		}

		response.OkResponse(w, http.StatusOK, SendResponse{Message: "sent"})
	}
}

func newHash() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

package send

import (
	"errors"
	"net/mail"
)

type SendRequest struct {
	Email string `json:"email"`
}

func (r SendRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	_, err := mail.ParseAddress(r.Email)
	return err
}

type SendResponse struct {
	Message string `json:"message"`
}

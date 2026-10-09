package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func DecodeAndValidate[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var req T
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req); err != nil {
		return nil, errors.New("invalid JSON body")
	}

	if err := validate.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			msgs := make([]string, 0, len(ve))
			for _, fe := range ve {
				msgs = append(msgs, fmt.Sprintf("%s: %s", fe.Field(), fe.Tag()))
			}
			return nil, errors.New(strings.Join(msgs, "; "))
		}
		return nil, errors.New("invalid request")
	}

	return &req, nil
}

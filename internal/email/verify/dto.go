package verify

type VerifyResponse struct {
	Message string `json:"message"`
	Email   string `json:"email,omitempty"`
}

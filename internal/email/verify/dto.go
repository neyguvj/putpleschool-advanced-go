package verify

type VerifyResponse struct {
	Verified bool   `json:"verified"`
	Email    string `json:"email,omitempty"`
}

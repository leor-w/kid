package apple

import "net/http"

const (
	VALIDATION_URL = "https://appleid.apple.com/auth/token"
)

type Client struct {
	validationURL string
	revokeURL     string
	client        *http.Client
}

func NewClient()

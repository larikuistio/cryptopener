package cryptopener

import (
	"net/http"
	"time"
)

// just a base struct
type cryptopener struct {
	client http.Client
}

func setupClient() http.Client {
	return http.Client{
		Timeout: time.Duration(60),
	}
}

func NewCryptopener() *cryptopener {
	cryptopener := cryptopener{
		client: setupClient(),
	}
	return &cryptopener
}

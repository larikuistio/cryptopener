package cryptopener

import (
	client "github.com/larikuistio/cryptopener/client"
)

// just a base struct
type cryptopener struct {
	client client.Client
}


func NewCryptopener(address string, entry string) *cryptopener {
	cryptopener := cryptopener{
		client: client.Client{
			Addr: address,
			Entrypoint: entry,
		},
	}
	return &cryptopener
}


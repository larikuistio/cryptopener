package cryptopener

import (
	client "github.com/larikuistio/cryptopener/client"
)

// Cryptopener a struct that defines cryptopener
type Cryptopener struct {
	// client for sending data to target server
	client client.Client
	// token mutator that creates new payloads
	mutator *TokenMutator
	ResultToken []string
}

// NewCryptopener creates new instance of Cryptopener
func NewCryptopener(address string, entry string) *Cryptopener {
	cryptopener := Cryptopener{
		client: client.Client{
			Addr: address,
			Entrypoint: entry,
		},
		mutator: NewMutator(),
	}
	return &cryptopener
}

func (p *Cryptopener) sendPayload(channel chan []byte, payload []byte) {
	channel <- []byte{}
	return
}

func (p *Cryptopener) analyseResponse(response []byte) {
	return
}

// Run starts BREACH attack
func (p *Cryptopener) Run() {
	channel := make(chan []byte, 1)
	for {
		// create new payload
		payload, _ := p.mutator.NewPayload(false)

		// send payload into a socket
		go p.sendPayload(channel, payload)
		defer p.analyseResponse(<-channel)
	}
}

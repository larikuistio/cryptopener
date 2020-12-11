package cryptopener

import (
	"log"
	"fmt"
	"net/http"
)


// Token mutator creates new token mutations
type TokenMutator struct {
	Tokens []byte
	target string
	port int
}

func NewMutator(target string, port int) *TokenMutator {
	return &TokenMutator{
		// for now make empty slice for used mutations
		Tokens: []byte{},
		target: target,
		port: port,
	}
}

// Create new payload
func (mutator *TokenMutator) newPayload() ([]byte, error) {
	return []byte{}, nil
}

func (mutator TokenMutator) SendNextPayload() {
	payload, err := mutator.newPayload()
	if err != nil {
		log.Fatalln("Could not create new payload")
		return
	}
	// hack together new url
	response, err != http.Get(fmt.Sprintf("%s%s", mutator.target, payload))
	if err != nil {

	}
}
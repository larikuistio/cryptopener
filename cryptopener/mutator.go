package cryptopener

import (
	"log"
	"fmt"
	"net/http"
	"unsafe"
)

const TOKENS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"


type tokenIterator struct {
	// list of all tokens
	tokens []byte
	previous int
}

func (iterator *tokenIterator) nextToken() byte {
	next_index := iterator.previous + 1
	if next_index > len(iterator.tokens) {
		next_index = 0
	}
	next := iterator.tokens[next_index]
	iterator.previous = next_index
	return next
}

// Token mutator creates new token mutations
type TokenMutator struct {
	target string
	port int
	previous []byte
	iterator tokenIterator
	result []byte
}

func NewMutator(target string, port int) *TokenMutator {
	return &TokenMutator{
		target: target,
		port: port,
		previous: []byte{},
		iterator: tokenIterator{
			tokens: []byte(TOKENS),
			previous: 0,
		},
	}
}

// Create new payload
func (mutator *TokenMutator) newPayload() ([]byte, error) {
	// new payload with
	new_payload := make([]byte, len(mutator.previous) + 1)

	// copy previous payload to new slice
	copy(new_payload, mutator.previous)
	new_payload = append(new_payload, mutator.iterator.nextToken())
	mutator.previous = new_payload
	return new_payload, nil
}

func analyseResponse(channel chan *http.Response) {
	response := <- channel
	// NOTE: this is absolute nono, but for now this have to do
	// so we are reading file size from memory
	size := unsafe.Sizeof(response)
	log.Println("Payload size:", size)
	return
}

// SendNextPayload generate and send new payload
func (mutator *TokenMutator) SendNextPayload() {
	payload, err := mutator.newPayload()
	if err != nil {
		log.Fatalln("Could not create new payload")
		return
	}
	// hack together new url
	response, err := http.Get(fmt.Sprintf("%s%s", mutator.target, payload))
	if err != nil {
		log.Fatalln("Could not send request")
		return
	}
	defer response.Body.Close()

	// create new channel for communication and send http response to channel
	channel := make(chan *http.Response)
	channel <- response

	// analyse response
	go analyseResponse(channel)
}

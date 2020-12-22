package cryptopener

import (
	"fmt"
	"log"
	"math"
)

const TOKENS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"


func createNewPayload(previousPayload []byte, length int) []byte {
	newPayload := make([]byte, length)
	if len(previousPayload) == 0 {
		return newPayload
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
		index: 0,
		previousPayload: []byte{},
		result: []string{},
		tokenCount: len(tokens),
	}
}

func (mutator TokenMutator) nextToken() byte {
	// increase index
	mutator.index ++
	nextToken := tokens[int(mutator.index % mutator.tokenCount)]
	return nextToken
}

// countPayloadLength counts length of next payload
func (mutator TokenMutator) countPayloadLength() (int, error) {
	if mutator.index == 0 {
		return 0, fmt.Errorf("Tried to call payload length count with 0")
	}
	return int(math.Floor(float64(mutator.index % mutator.tokenCount))), nil
}

// Creates new payload
func (mutator *TokenMutator) newPayload(savePrevious bool) ([]byte, error) {
	nextPayloadLength, err := mutator.countPayloadLength()
	if err != nil {
		log.Printf("Could not count payload length, errpr %e", err)
		return nil, err
	}
	newPayload := createNewPayload(mutator.previousPayload, nextPayloadLength)
	nextToken := mutator.nextToken()
	if nextPayloadLength > len(newPayload) || savePrevious {
		if savePrevious {
			mutator.result = append(mutator.result, string(mutator.previousPayload ))
		}
		newPayload = append(newPayload, nextToken)
		mutator.previousPayload = newPayload
		return newPayload, nil
	}
	// remove last element from slice
	newPayload = newPayload[:len(newPayload) - 1]
	newPayload = append(newPayload, nextToken)
	return newPayload, nil
}

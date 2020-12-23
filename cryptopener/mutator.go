package cryptopener

import (
	"fmt"
	"log"
	"math"
)

// tokens is a array of used tokens
const tokens = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"


func createNewPayload(previousPayload []byte, length int) []byte {
	newPayload := make([]byte, length)
	if len(previousPayload) == 0 {
		return newPayload
	}
	copy(newPayload, previousPayload)
	return newPayload
}

// TokenMutator creates new token mutations
type TokenMutator struct {
	// previously used payload
	previousPayload []byte
	// correctly quessed tokens
	result []string
	index, tokenCount int
}

// NewMutator creates a new payload mutator
func NewMutator() *TokenMutator {
	return &TokenMutator{
		index: 0,
		previousPayload: []byte{},
		result: []string{},
		tokenCount: len(tokens),
	}
}

func (mutator TokenMutator) nextToken() byte {
	nextToken := tokens[int(mutator.index % mutator.tokenCount)]
	mutator.index++
	return nextToken
}

// countPayloadLength counts length of next payload
func (mutator TokenMutator) countPayloadLength() (int, error) {
	if mutator.index == 0 {
		return 0, fmt.Errorf("Tried to call payload length count with 0")
	}
	return int(math.Floor(float64(mutator.index / mutator.tokenCount))), nil
}

// NewPayload creates new payload
func (mutator *TokenMutator) NewPayload(savePrevious bool) ([]byte, error) {
	nextToken := mutator.nextToken()
	nextPayloadLength, err := mutator.countPayloadLength()
	if err != nil {
		log.Printf("Could not count payload length, errpr %e", err)
		return nil, err
	}
	newPayload := createNewPayload(mutator.previousPayload, nextPayloadLength)
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

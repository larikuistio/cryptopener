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
	previousPayload , result[]byte
	index, tokenCount int
}

// NewMutator creates a new payload mutator
func NewMutator() *TokenMutator {
	return &TokenMutator{
		index: 0,
		previousPayload: []byte{},
		result: []byte{},
		tokenCount: len(tokens),
	}
}

func (mutator *TokenMutator) NextToken() byte {
	nextToken := tokens[int(mutator.index % mutator.tokenCount)]
	mutator.index++
	return nextToken
}

// countPayloadLength counts length of next payload
func (mutator TokenMutator) countPayloadLength() (int, error) {
	if mutator.index == 0 {
		return 0, fmt.Errorf("Tried to call payload length count with 0")
	}
	var length int
	if mutator.index % mutator.tokenCount == 0 {
		length = int(math.Floor(float64(mutator.index / mutator.tokenCount)))
	} else {
		length = int(math.Floor(float64(mutator.index / mutator.tokenCount))) + 1
	}
	return length, nil
}

// NewPayload creates new payload
func (mutator *TokenMutator) NewPayload(savePrevious bool) ([]byte, error) {
	nextToken := mutator.NextToken()
	nextPayloadLength, err := mutator.countPayloadLength()
	if err != nil {
		log.Printf("Could not count payload length, errpr %e", err)
		return nil, err
	}

	var newPayload []byte
	if len(mutator.result) != 0 {
		newPayload = createNewPayload(mutator.result, nextPayloadLength)
	} else {
		newPayload = createNewPayload(mutator.previousPayload, nextPayloadLength)
	}

	if nextPayloadLength > len(mutator.previousPayload) || len(newPayload) == 0 || savePrevious {
		log.Printf("old new %s %d %d", newPayload, nextPayloadLength, len(mutator.previousPayload))
		if savePrevious {
			mutator.result = append(mutator.result, mutator.previousPayload...)
		}
		newPayload = append(newPayload, nextToken)
		mutator.previousPayload = newPayload
		return newPayload, nil
	}
	newPayload = newPayload[:len(newPayload) - 1]
	newPayload = append(newPayload, nextToken)
	mutator.previousPayload = newPayload
	return newPayload, nil
}

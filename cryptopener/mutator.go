package cryptopener

import (
	"bytes"
	"math"
	"fmt"
)

// tokens is a array of used tokens
const tokens = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"


func createNewPayload(previuosPayload []byte, length int) []byte {
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
func NewMutator(target string, port int) *TokenMutator {
	return &TokenMutator{
		index: 0,
		previousPayload: []byte{},
		result: []string{}
		tokenCount: len(tokens)
		iterator: tokenIterator{
			tokens: []byte(tokens),
			previous: 0,
		},
	}
}

func (mutator TokenMutator) nextToken() byte {
	// increase index
	mutator.index += 1
	nextToken := mutator.tokens[int(mutator.index % mutator.tokenCount)]
	return nextToken
}

// countPayloadLength counts length of next payload
func (mutator TokenMutator) countPayloadLength() (int, error) {
	if mutator.index == 0 {
		return 0, fmt.Error("Tried to call payload length count with 0")
	}
	return int(math.Floor(mutator.index % mutator.tokenCount)), nil
}

// Creates new payload
func (mutator *TokenMutator) newPayload(savePrevious bool) ([]byte, error) {
	nextPayloadLength := mutator.countPayloadLength()
	newPayload := createNewPayload(mutator.previuosPayload, nextPayloadLength)
	if nextPayloadLength > len(newPayload) || savePrevious {
		if savePrevious {
			mutator.result = append(mutator.result, string(mutator.previuosPayload))
		}
		newPayload = append(newPayload, nextToken)
		mutator.previous = newPayload
		return newPayload, nil
	}
	// remove last element from slice
	newPayload = newPayload[:len(newPayload) - 1]
	newPayload = append(newPayload, nextToken)
	return newPayload, nil
}

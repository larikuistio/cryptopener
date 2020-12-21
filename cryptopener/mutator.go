package cryptopener

import "bytes"

// tokens is a array of used tokens
const tokens = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"


// TokenMutator creates new token mutations
type TokenMutator struct {
	// previously used payload
	previousPayload []byte
	// correctly quessed tokens
	result []string
	// position in the token array
	index int
	payloadLength int
}

// NewMutator creates a new payload mutator
func NewMutator(target string, port int) *TokenMutator {
	return &TokenMutator{
		index: 0,
		previousPayload: []byte{},
		result: []string{}
		payloadLength: 1
		iterator: tokenIterator{
			tokens: []byte(tokens),
			previous: 0,
		},
	}
}

func (mutator TokenMutator) nextToken() byte {
	// increase index
	mutator.index += 1
	remainder := mutator.index % len(tokens)
	if remainder == 0{
		mutator.payloadLength += 1
	}
	nextToken := mutator.tokens[remainder]
	return nextToken
}

// Creates new payload
func (mutator *TokenMutator) newPayload(savePrevious bool) ([]byte, error) {
	nextToken := mutator.nextToken()
	newPayload := make([]byte, mutator.payloadLength)
	if mutator.payloadLength > len(mutator.previousPayload) || savePrevious {
		// copy previous payload to new slice
		copy(newPayload, mutator.previousPayload)
		newPayload = append(newPayload, nextToken)
		mutator.previous = newPayload
		return newPayload, nil
	}
	// FIXME: clean up this mess
	newPayload = append(newPayload, mutator.previousPayload[:len(mutator.previousPayload - 1)]...)
	newPayload = append(newPayload, nextToken)
	return newPayload, nil
}

package cryptopener

import "log"

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

type guessToken struct {
	index int
}

func (guess *guessToken) nextToken(tokenCount int) byte {
	nextToken := tokens[int(guess.index % tokenCount)]
	guess.index++
	return nextToken
}

func (guess guessToken) getCurrentToken() byte {
 	return tokens[guess.index]
}

func (guess guessToken) movePosition(tokenCount int) bool {
	remainder := int(guess.index % tokenCount)
	if remainder == 0 && guess.index != 0 {
		return true
	}
	return false
}

// TokenMutator creates new token mutations
type TokenMutator struct {
	// previously used payload
	previousPayload , result[]byte
	previousIndex, tokenCount, position int
	tokenMap map[int]guessToken
}

// NewMutator creates a new payload mutator
func NewMutator() *TokenMutator {
	mutator := TokenMutator{
		previousIndex: 0,
		position: 0,
		previousPayload: []byte{},
		result: []byte{},
		tokenCount: len(tokens),
		tokenMap: make(map[int]guessToken),
	}
	// create initial token
	mutator.tokenMap[0] = guessToken{
		index: 0,
	}
	return &mutator
}

func (mutator *TokenMutator) addToken() {
	newToken := guessToken{
		index: 0,
	}
	mutator.tokenMap[len(mutator.tokenMap)] = newToken
	mutator.position = 0
}

// NewPayload creates new payload
func (mutator *TokenMutator) NewPayload(savePrevious bool) ([]byte, error) {
	if mutator.position >= len(mutator.tokenMap) && mutator.tokenMap[mutator.position].movePosition(mutator.tokenCount) {
		mutator.addToken()
	}

	var newPayload []byte
	for pos, elem := range mutator.tokenMap {
		if pos == mutator.position {
			log.Panicf("HERE %d %d %s", pos, mutator.position, string(elem.getCurrentToken()))
			newPayload = append(newPayload, elem.nextToken(mutator.tokenCount))
			if len(mutator.tokenMap) > 1 {
				mutator.position++
			}
		} else {
			newPayload = append(newPayload, elem.getCurrentToken())
		}
	}
	return newPayload, nil
}

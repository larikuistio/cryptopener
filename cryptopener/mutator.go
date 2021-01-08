package cryptopener

import (
	"fmt"
	"log"
	//"log"
	"math/big"
)

// tokens is a array of used tokens
const tokens = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Factorial counts factorial for value
func factorial(val int) (*big.Int, error) {
	if val <= 0 {
		return &big.Int{}, fmt.Errorf("Cannot create factorial from negative value")
	}
	result := big.NewInt(1)
	for i := 1; i <= val; i++ {
		j := big.NewInt(int64(i))
		result.Mul(result, j)

	}
	return result, nil
}

// CountPermutations counts amount of permutations for set of values
func countPermutations(totalValues int, itemCount int) (int64, error) {
	totalPermutations, err := factorial(totalValues)
	if err != nil {
		return int64(-1), err
	}

	divider := totalValues - itemCount
	if divider <= 0 {
		return int64(-1), fmt.Errorf("Cannot create factorial from negative value")
	}

	subPermutations, _ := factorial(divider)

	result := new(big.Int).Div(totalPermutations, subPermutations).Uint64()
	return int64(result), nil
}

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

func (guess guessToken) getCurrentToken(tokenCount int) byte {
 	return tokens[int(guess.index % tokenCount)]
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
	tokenCount, position, lastIndexPos int
	maxPermutations, mutations int64
	tokenMap map[int]*guessToken
}

// NewMutator creates a new payload mutator
func NewMutator() *TokenMutator {
	mutator := TokenMutator{
		mutations: 0,
		position: 0,
		maxPermutations: int64(len(tokens)),
		previousPayload: []byte{},
		result: []byte{},
		tokenCount: len(tokens),
		tokenMap: make(map[int]*guessToken),
		lastIndexPos: 0,
	}
	// create initial token
	mutator.tokenMap[0] = &guessToken{
		index: 0,
	}
	return &mutator
}

func (mutator *TokenMutator) addToken() {
	newToken := &guessToken{
		index: 0,
	}
	mutator.lastIndexPos++
	mutator.tokenMap[mutator.lastIndexPos] = newToken
	mutator.position = 0
}

// NewPayload creates new payload
func (mutator *TokenMutator) NewPayload(savePrevious bool) ([]byte, error) {
	// FIXME: is heavy please cache
	if int(mutator.tokenMap[mutator.lastIndexPos].index % mutator.tokenCount) == 0 && mutator.tokenMap[mutator.lastIndexPos].index != 0{
		mutator.addToken()
	}

	var newPayload []byte
	for pos := 0; pos <= len(mutator.tokenMap) - 1; pos++ {
		elem := mutator.tokenMap[pos]
		if pos == mutator.position {
			newPayload = append(newPayload, elem.nextToken(mutator.tokenCount))
			if len(mutator.tokenMap) > 1 {
				mutator.position++
			}
		} else {
			newPayload = append(newPayload, elem.getCurrentToken(mutator.tokenCount))
		}
	}
	mutator.mutations++
	return newPayload, nil
}

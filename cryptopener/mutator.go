package cryptopener


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
	previousPayload, result[]byte
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
}

// NewPayload creates new payload
func (mutator *TokenMutator) NewPayload(savePrevious bool) ([]byte, error) {
	if savePrevious {
		mutator.result = append(mutator.result, mutator.previousPayload...)
		// clear existing token map and create new one
		mutator.tokenMap = make(map[int]*guessToken)
		mutator.tokenMap[0] = &guessToken{
			index: 0,
		}

	}
	// checks if all permutations are done
	var sum int
	for _, i := range mutator.tokenMap {
		sum += int(i.index % mutator.tokenCount)
	}

	if sum == 0 && mutator.mutations >= 62 {
		mutator.addToken()
		mutator.position = 0
	} else if mutator.position > len(mutator.tokenMap) {
		mutator.position = 0
	}

	var newPayload []byte
	if len(mutator.result) > 0 {
		newPayload = append(newPayload, mutator.result...)
	}

	for pos := 0; pos <= len(mutator.tokenMap) - 1; pos++ {
		elem := mutator.tokenMap[pos]
		if pos <= mutator.position {
			newPayload = append(newPayload, elem.nextToken(mutator.tokenCount))
		} else {
			newPayload = append(newPayload, elem.getCurrentToken(mutator.tokenCount))
		}
	}

	// keep track of previous payload
	mutator.previousPayload = newPayload

	// move to next token
	mutator.position++
	mutator.mutations++
	return newPayload, nil
}

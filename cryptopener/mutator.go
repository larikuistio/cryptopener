package cryptopener


const TOKENS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"


type tokenIterator struct {
	// list of all tokens
	tokens []byte
	previous int
}

func (iterator *tokenIterator) nextToken() byte {
	next_index := iterator.previous + 1
	if next_index > len(iterator.tokens) {
		next_index = 0
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
		target: target,
		port: port,
		previous: []byte{},
		iterator: tokenIterator{
			tokens: []byte(TOKENS),
			previous: 0,
		},
	}
}

// Create new payload
func (mutator *TokenMutator) newPayload() ([]byte, error) {
	// new payload with
	new_payload := make([]byte, len(mutator.previous) + 1)

	// copy previous payload to new slice
	copy(new_payload, mutator.previous)
	new_payload = append(new_payload, mutator.iterator.nextToken())
	mutator.previous = new_payload
	return new_payload, nil
}

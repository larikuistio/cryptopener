package cryptopener


// Token mutator creates new token mutations
type TokenMutator struct {
	tokens []byte
}

func NewMutator () *TokenMutator {
	return &TokenMutator{
		// for now make empty slice for used mutations
		tokens: []byte{},
	}
}

// Create new payload
func (mutator *TokenMutator) NewPayload() *[]byte {
	return &[]byte{}
}
package models

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Integer type Number for GraphQL indexes and counters.
type Number hexutil.Big

// Map this custom Go type to GraphQL scalar type on the schema.
func (Number) ImplementsGraphQLType(name string) bool {
	return "Number" == name
}

// MarshalJSON implements the json.Marshaler interface.
func (n Number) MarshalJSON() ([]byte, error) {
	b := hexutil.Big(n)
	return []byte(b.ToInt().String()), nil
}

// Decode the Number into the Go type.
//
// This will be called whenever you use Number scalar in input.
func (n *Number) UnmarshalGraphQL(input interface{}) error {
	return n.UnmarshalGraphQL(input)
}

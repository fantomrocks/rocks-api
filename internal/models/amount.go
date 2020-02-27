package models

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
	"math/big"
)

// Defines Amount entity for GraphQL
type Amount struct {
	decimal.Decimal
}

// Map this custom Go structure to GraphQL scalar type on the schema.
func (Amount) ImplementsGraphQLType(name string) bool {
	return "Amount" == name
}

// Decode the Amount into the Go structure.
// This will be called whenever you use Amount scalar in input.
func (a *Amount) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {
	case string:
		return a.UnmarshalJSON([]byte(input))
	default:
		return fmt.Errorf("wrong type")
	}
}

// Convert fraction amount to token amount.
func (a *Amount) ToFTM() *Amount {
	return &Amount{Decimal: a.Decimal.Div(decimal.New(1, 18))}
}

// Convert Amount to HEX value appropriate for tokens transfer
// Warning, only integer part of the value is considered! We need to make sure this will work as intended.
func (a *Amount) ToHex() string {
	return hexutil.EncodeBig(big.NewInt(a.IntPart()))
}

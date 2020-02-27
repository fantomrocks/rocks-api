package models

import "github.com/graph-gophers/graphql-go"

// Define a Blockchain Transaction entity.
type BcBlock struct {
	Hash      string
	Number    Number
	TimeStamp graphql.Time
	TxHashes  []string
}

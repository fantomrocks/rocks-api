package models

import "github.com/graph-gophers/graphql-go"

// Define Transaction entity.
type Transaction struct {
	Id          string
	FromAccount *Account
	ToAccount   *Account
	Amount      *Amount
	TimeStamp   *graphql.Time
}

package types

import (
	"fantomrocks-api/internal/models"
	"fantomrocks-api/internal/repository"
	"github.com/graph-gophers/graphql-go"
)

// Define Transaction type.
type Transaction struct {
	repo *repository.Repository
	tr   *models.Transaction
}

// Make new Transaction.
func NewTransaction(tr *models.Transaction, repo *repository.Repository) *Transaction {
	return &Transaction{
		repo: repo,
		tr:   tr,
	}
}

// Properly resolve the GraphQL.ID where needed.
func (t *Transaction) ID() graphql.ID {
	return graphql.ID(t.tr.Id)
}

// Resolve source Account
func (t *Transaction) From() *Account {
	return NewAccount(t.tr.FromAccount, t.repo)
}

// Resolve destination Account
func (t *Transaction) To() *Account {
	return NewAccount(t.tr.ToAccount, t.repo)
}

// Resolve Transaction Amount
func (t *Transaction) Amount() models.Amount {
	return *t.tr.Amount
}

// Resolve Transaction time stamp
func (t *Transaction) TimeStamp() graphql.Time {
	return *t.tr.TimeStamp
}

package types

import (
	"fantomrocks-api/internal/models"
	"fantomrocks-api/internal/repository"
	"github.com/graph-gophers/graphql-go"
	"strconv"
)

// Define Account Resolver for GraphQL
type Account struct {
	repo *repository.Repository
	acc  *models.Account
}

// Make new Account
func NewAccount(act *models.Account, repo *repository.Repository) *Account {
	return &Account{
		repo: repo,
		acc:  act,
	}
}

// Properly resolve the GraphQL.ID where needed.
// GraphQL.ID is of type <string> so we need a conversion from internal int64.
func (a *Account) ID() graphql.ID {
	return graphql.ID(strconv.FormatInt(a.acc.Id, 10))
}

// Get Account balance for GraphQL.
func (a *Account) Balance() (models.Amount, error) {
	balance, err := a.repo.Rpc.AccountBalance(a.acc.Address)
	return *balance, err
}

// Resolve Account name
func (a *Account) Name() string {
	return a.acc.Name
}

// Resolve Account address
func (a *Account) Address() string {
	return a.acc.Address
}

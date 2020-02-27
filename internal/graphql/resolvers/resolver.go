package resolvers

import (
	"fantomrocks-api/internal/graphql/inputs"
	"fantomrocks-api/internal/graphql/types"
	"fantomrocks-api/internal/models"
	"fantomrocks-api/internal/repository"
	"fantomrocks-api/internal/services"
	"github.com/graph-gophers/graphql-go"
)

// Define new Use Cases interface
type UseCases interface {
	// Query fo Accounts
	Account(args *struct{ Id *graphql.ID }) (*types.Account, error)
	RandomAccount() (*types.Account, error)
	Accounts(*struct{ List *[]graphql.ID }) ([]*types.Account, error)

	// Query for Pairs
	Pair() (*types.AccountPair, error)
	Pairs() ([]*types.AccountPair, error)

	// Query for Transactions and Blocks
	BlockchainTransaction(*struct{ Hash graphql.ID }) (*types.BlockchainTransaction, error)

	// Mutation
	Transfer(*struct{ ToTransfer inputs.TransferInput }) (*types.Transaction, error)
	Burst(*struct {
		FromAccountId graphql.ID
		Amount        models.Amount
		TargetsCount  int32
	}) ([]*types.Transaction, error)
}

// Defines root resolver to be used to define entry points.
type Resolver struct {
	log services.Logger
	*repository.Repository
}

// Create new
func NewResolver(repo *repository.Repository, log services.Logger) UseCases {
	return &Resolver{log: log, Repository: repo}
}

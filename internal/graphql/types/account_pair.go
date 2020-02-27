package types

import (
	"fantomrocks-api/internal/models"
	"fantomrocks-api/internal/repository"
)

// Define Account Pair type.
type AccountPair struct {
	repo *repository.Repository
	pair *models.AccountPair
}

// Make new Account Pair
func NewAccountPair(pair *models.AccountPair, repo *repository.Repository) *AccountPair {
	return &AccountPair{
		repo: repo,
		pair: pair,
	}
}

// Resolve the first Account of the Pair
func (ap *AccountPair) One() *Account {
	return NewAccount(ap.pair.One, ap.repo)
}

// Resolve the second Account of the Pair
func (ap *AccountPair) Two() *Account {
	return NewAccount(ap.pair.Two, ap.repo)
}

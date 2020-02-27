package resolvers

import (
	"fantomrocks-api/internal/graphql/types"
	"github.com/graph-gophers/graphql-go"
	"strconv"
)

// Implements Query.account for an Account specified by ID or random
func (rs *Resolver) Account(args *struct{ Id *graphql.ID }) (*types.Account, error) {
	// random account requested?
	if args == nil || args.Id == nil {
		return rs.RandomAccount()
	}

	// get the account id
	id, err := strconv.Atoi(string(*args.Id))
	if err != nil {
		rs.log.Errorf("GQL->Mutation->Account(): Invalid account ID [%s]. %s", args.Id, err.Error())
		return nil, err
	}

	acc, err := rs.Db.AccountById(id)
	if err != nil {
		rs.log.Errorf("GQL->Query->Account(): Can not get Account. %s", err.Error())
		return nil, err
	}

	return types.NewAccount(acc, rs.Repository), nil
}

// Implements Query.account for a random Account
func (rs *Resolver) RandomAccount() (*types.Account, error) {
	// get random account info from DB
	acc, err := rs.Db.RandomAccount()
	if err != nil {
		rs.log.Errorf("GQL->Query->RandomAccount(): Can not get a random Account. %s", err.Error())
		return nil, err
	}

	return types.NewAccount(acc, rs.Repository), nil
}

// Implements Query.accounts GraphQL entry point
func (rs *Resolver) Accounts(args *struct{ List *[]graphql.ID }) ([]*types.Account, error) {
	// all accounts, or specified subset?
	if args == nil || args.List == nil {
		return rs.AllAccounts()
	}

	// log the action
	rs.log.Debugf("GQL->Query->Accounts(): List of specified accounts is constructed.")

	// loop requested accounts and construct the response
	result := make([]*types.Account, 0)
	for _, aid := range *args.List {
		// try to convert incoming account id
		id, err := strconv.Atoi(string(aid))
		if err != nil {
			rs.log.Debugf("GQL->Query->Accounts(): Invalid incoming account id '%s'.", aid)
			continue
		}

		// try to find the account in our database
		a, err := rs.Db.AccountById(id)
		if err != nil {
			rs.log.Debugf("GQL->Query->Accounts(): Incoming account '%s' not found.", aid)
			continue
		}

		// add the account to output
		result = append(result, types.NewAccount(a, rs.Repository))
	}

	return result, nil
}

// Implements Query.accounts for all the Account in the app
func (rs *Resolver) AllAccounts() ([]*types.Account, error) {
	// log the action
	rs.log.Debugf("GQL->Query->AllAccounts(): List of all accounts is constructed.")

	// get the list of accounts from the Account Service
	accounts, err := rs.Db.AllAccounts()
	if err != nil {
		rs.log.Errorf("GQL->Query->AllAccounts(): Can not get list of Accounts. %s", err.Error())
		return make([]*types.Account, 0), err
	}

	// prep resolvers for the accounts loaded
	result := make([]*types.Account, 0)
	for _, a := range accounts {
		result = append(result, types.NewAccount(a, rs.Repository))
	}

	return result, nil
}

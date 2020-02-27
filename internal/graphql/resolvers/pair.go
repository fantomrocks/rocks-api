package resolvers

import "fantomrocks-api/internal/graphql/types"

// Implements Query.pair GraphQL entry point for random Account Pair selection
func (rs *Resolver) Pair() (*types.AccountPair, error) {
	// log the action
	rs.log.Debugf("GQL->Query->Pair(): Random account pair is prepared.")

	// get the pair
	pair, err := rs.Db.RandomPair()

	if err != nil {
		rs.log.Errorf("GQL->Query->Pair(): Can not get random Account Pair. %s", err.Error())
		return &types.AccountPair{}, err
	}

	return types.NewAccountPair(pair, rs.Repository), nil
}

// Implements Query.pairs GraphQL entry point
func (rs *Resolver) Pairs() ([]*types.AccountPair, error) {
	// log the action
	rs.log.Debugf("GQL->Query->Pairs(): List of all account pairs is constructed.")

	// get the list of pairs from Account Service
	pairs, err := rs.Db.AllPairs()
	if err != nil {
		rs.log.Errorf("GQL->Query->AccountPairs(): Can not get list of Account Pairs. %s", err.Error())
		return make([]*types.AccountPair, 0), err
	}

	// prep resolvers for the Account Pairs loaded
	result := make([]*types.AccountPair, 0)
	for _, p := range pairs {
		result = append(result, types.NewAccountPair(p, rs.Repository))
	}

	return result, nil
}

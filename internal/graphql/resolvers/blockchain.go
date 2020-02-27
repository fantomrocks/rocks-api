package resolvers

import (
	"fantomrocks-api/internal/graphql/types"
	"github.com/graph-gophers/graphql-go"
)

// Get details of a blockchain Transaction by its identifier / hash
func (rs *Resolver) BlockchainTransaction(args *struct{ Hash graphql.ID }) (*types.BlockchainTransaction, error) {
	tx, err := rs.Rpc.TransactionByHash(string(args.Hash))
	if err != nil {
		rs.log.Errorf("GQL->Query->TransactionByHash(): Can not get Transaction. %s", err.Error())
		return nil, err
	}

	return types.NewBlockchainTransaction(tx, rs.Repository), nil
}

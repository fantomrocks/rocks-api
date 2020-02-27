package resolvers

import (
	"fantomrocks-api/internal/graphql/inputs"
	"fantomrocks-api/internal/graphql/types"
	"fantomrocks-api/internal/models"
	"fantomrocks-api/internal/repository/rpc"
	"github.com/graph-gophers/graphql-go"
	"strconv"
)

// Implements Mutation.transfer GraphQL entry point for sending a single transaction between a pair of internal accounts.
func (rs *Resolver) Transfer(args *struct{ ToTransfer inputs.TransferInput }) (*types.Transaction, error) {
	// get the source id
	fid, err := strconv.Atoi(string(args.ToTransfer.FromAccountId))
	if err != nil {
		// log the error and quit
		rs.log.Errorf("GQL->Mutation->Transfer(): Invalid source account ID [%s]. %s", args.ToTransfer.FromAccountId, err.Error())
		return nil, err
	}

	// get the source
	from, err := rs.Db.AccountById(fid)
	if err != nil {
		// log the error and quit
		rs.log.Errorf("GQL->Mutation->Transfer(): Source account not found for account id [%s]. %s", args.ToTransfer.FromAccountId, err.Error())
		return nil, err
	}

	// get the target id
	tid, err := strconv.Atoi(string(args.ToTransfer.ToAccountId))
	if err != nil {
		// log the error and quit
		rs.log.Errorf("GQL->Mutation->Transfer(): Invalid destination account ID [%s]. %s", args.ToTransfer.ToAccountId, err.Error())
		return nil, err
	}

	// get the source
	to, err := rs.Db.AccountById(tid)
	if err != nil {
		// log the error and quit
		rs.log.Errorf("GQL->Mutation->Transfer(): Destination account not found for account id [%s]. %s", args.ToTransfer.ToAccountId, err.Error())
		return nil, err
	}

	// log the action
	rs.log.Debugf("GQL->Mutation->Transfer(): Sending %s FTM tokens [%s -> %s].", args.ToTransfer.Amount.ToFTM(), from.Name, to.Name)

	// do the transfer
	tr, err := rs.Rpc.TransferTokens(from, to, args.ToTransfer.Amount)
	if err != nil {
		// log the action
		rs.log.Errorf("GQL->Mutation->Transfer(): Can not send tokens. %s", err.Error())
		return nil, err
	}

	// return nothing
	return types.NewTransaction(tr, rs.Repository), nil
}

// Implements Mutation.burst GraphQL entry point for sending blob of transactions to Opera node for parallel processing.
func (rs *Resolver) Burst(args *struct {
	FromAccountId graphql.ID
	Amount        models.Amount
	TargetsCount  int32
}) ([]*types.Transaction, error) {
	// prep empty result set slice
	result := make([]*types.Transaction, 0)

	// get the source account id
	id, err := strconv.Atoi(string(args.FromAccountId))
	if err != nil {
		rs.log.Errorf("GQL->Mutation->Burst(): Invalid source account ID [%s]. %s", args.FromAccountId, err.Error())
		return result, err
	}

	// get the source account details
	from, err := rs.Db.AccountById(id)
	if err != nil {
		rs.log.Errorf("GQL->Mutation->Burst(): Source account not found for account id [%s]. %s", args.FromAccountId, err.Error())
		return nil, err
	}

	// get list of random account to work with on the burst
	accounts, err := rs.Db.RandomAccounts(int(args.TargetsCount), []*models.Account{from})
	if err != nil {
		// log the error and quit
		rs.log.Errorf("GQL->Mutation->Burst(): Could not get a list of target accounts. %s", err.Error())
		return result, err
	}

	// do we have any accounts to process?
	if 0 == len(accounts) {
		rs.log.Errorf("GQL->Mutation->Burst(): Could not continue, requested %d but no accounts found.", args.TargetsCount)
		return result, nil
	}

	// prep channel to receive transactions
	trs := make(chan *models.Transaction)
	defer close(trs)

	// inform
	rs.log.Debugf("GQL->Mutation->Burst(): Sending %d transactions.", len(accounts))

	// start sending in parallel
	for _, account := range accounts {
		// do actual sending
		go func(acc *models.Account, chain rpc.BlockChain) {
			// try to push the transfer
			tr, err := chain.TransferTokens(from, acc, args.Amount)
			if err != nil {
				rs.log.Errorf("GQL->Mutation->Burst(): Can not send tokens from %s to %s. %s", from.Name, acc.Name, err.Error())
			}

			// send the transaction to channel (or nil if the transaction failed)
			trs <- tr
		}(account, rs.Repository.Rpc)
	}

	// we know exactly how many we should get; extract from channel and prep valid TXes for output
	for i := 0; i < len(accounts); i++ {
		tr := <-trs
		if tr != nil {
			result = append(result, types.NewTransaction(tr, rs.Repository))
		}
	}

	// inform
	rs.log.Debugf("GQL->Mutation->Burst(): Done [%d].", len(accounts))

	// return what we've got here
	return result, nil
}

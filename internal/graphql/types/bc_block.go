package types

import (
	"fantomrocks-api/internal/models"
	"fantomrocks-api/internal/repository"
	"github.com/graph-gophers/graphql-go"
)

// Define Blockchain Block Resolver for GraphQL
type BlockchainBlock struct {
	repo *repository.Repository
	blk  *models.BcBlock
}

// Make new Account
func NewBlockchainBlock(blk *models.BcBlock, repo *repository.Repository) *BlockchainBlock {
	return &BlockchainBlock{
		repo: repo,
		blk:  blk,
	}
}

// Properly resolve the GraphQL.ID where needed.
func (b *BlockchainBlock) Hash() graphql.ID {
	return graphql.ID(b.blk.Hash)
}

// Number of the Block in the Blockchain.
func (b *BlockchainBlock) Number() models.Number {
	return b.blk.Number
}

// Resolve Transaction time stamp.
func (b *BlockchainBlock) TimeStamp() graphql.Time {
	return b.blk.TimeStamp
}

// Resolve list of TX hashes.
func (b *BlockchainBlock) TxHashes() []string {
	return b.blk.TxHashes
}

// Resolve full list of Transactions in the Block.
func (b *BlockchainBlock) Transactions() []*BlockchainTransaction {
	tx := make([]*BlockchainTransaction, 0)

	// any hashes registered in the block?
	if 0 < len(b.blk.TxHashes) {
		// loop all hashes and get corresponding transactions
		for _, h := range b.blk.TxHashes {
			t, err := b.repo.Rpc.TransactionByHash(h)
			if err == nil {
				// add the transaction detail to the output
				tx = append(tx, NewBlockchainTransaction(t, b.repo))
			} else {
				// log the error
				b.repo.Log.Debugf("GQL->BlockchainBlock(): Could not resolve transaction; %s", err.Error())
			}
		}
	}

	return tx
}

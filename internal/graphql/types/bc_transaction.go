package types

import (
	"fantomrocks-api/internal/models"
	"fantomrocks-api/internal/repository"
	"github.com/graph-gophers/graphql-go"
)

// Define Blockchain Transaction Resolver for GraphQL
type BlockchainTransaction struct {
	repo *repository.Repository
	tx   *models.BcTransaction
}

// Make new Account
func NewBlockchainTransaction(tx *models.BcTransaction, repo *repository.Repository) *BlockchainTransaction {
	return &BlockchainTransaction{
		repo: repo,
		tx:   tx,
	}
}

// Properly resolve the GraphQL.ID where needed.
func (t *BlockchainTransaction) Hash() graphql.ID {
	return graphql.ID(t.tx.Hash)
}

// Resolve the senders address.
func (t *BlockchainTransaction) From() string {
	return t.tx.From
}

// Resolve the receivers address.
func (t *BlockchainTransaction) To() *string {
	return t.tx.To
}

// Resolve the value transferred.
func (t *BlockchainTransaction) Value() models.Amount {
	return t.tx.Value
}

// Resolve the data attached to the transaction on sending.
func (t *BlockchainTransaction) Input() string {
	return t.tx.Input
}

// Resolve the number of transactions sender sent prior this one.
func (t *BlockchainTransaction) Nonce() int32 {
	return int32(t.tx.Nonce)
}

// Resolve the index of the Transaction in its block.
func (t *BlockchainTransaction) TxIndex() *int32 {
	return t.tx.TxIndex
}

// Resolve the transaction gas price.
func (t *BlockchainTransaction) GasPrice() models.Amount {
	return t.tx.GasPrice
}

// Resolve the amount of gas used by the transaction.
func (t *BlockchainTransaction) GasUsed() models.Amount {
	return t.tx.GasUsed
}

// Resolve the maximal amount of gas provided by sender for the transaction.
func (t *BlockchainTransaction) GasLimit() models.Amount {
	return t.tx.GasLimit
}

// Resolve the transaction fee paid by the sender.
func (t *BlockchainTransaction) Fee() models.Amount {
	return t.tx.Fee
}

// Resolve the Block this transaction belongs to.
func (t *BlockchainTransaction) Block() *BlockchainBlock {
	// just return no-block
	if nil == t.tx.BlockHash {
		return nil
	}

	b, err := t.repo.Rpc.BlockByHash(*t.tx.BlockHash)
	if err != nil {
		t.repo.Log.Errorf("GQL->BlockchainTransaction():: Block not loaded! %s", err)
		return nil
	}

	return NewBlockchainBlock(b, t.repo)
}

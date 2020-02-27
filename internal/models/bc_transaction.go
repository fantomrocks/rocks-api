package models

// Define a Blockchain Transaction entity.
type BcTransaction struct {
	Hash      string
	From      string
	To        *string
	Value     Amount
	Input     string
	Nonce     uint
	GasLimit  Amount
	GasUsed   Amount
	GasPrice  Amount
	Fee       Amount
	TxIndex   *int32
	BlockHash *string
}

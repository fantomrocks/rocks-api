package rpc

import (
	"fantomrocks-api/internal/models"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/graph-gophers/graphql-go"
	"time"
)

// Get a raw Transaction information for given tx hash.
func (rpc *Rpc) BlockByHash(hash string) (*models.BcBlock, error) {
	// unlock the source account
	rpc.log.Debugf("RPC->BlockByHash(): Loading block details for [%s]", hash)

	// container for raw data
	var raw struct {
		Hash         string       `json:"hash"`
		Number       hexutil.Big  `json:"number"`
		Miner        string       `json:"miner"`
		GasLimit     hexutil.Big  `json:"gasLimit"`
		GasUsed      hexutil.Big  `json:"gasUsed"`
		Timestamp    hexutil.Uint `json:"timestamp"`
		Transactions []string     `json:"transactions"`
	}

	// call for data
	err := rpc.Call(&raw, "eth_getBlockByHash", hash, false)
	if err != nil {
		rpc.log.Errorf("RPC->BlockByHash(): Error! %s", err.Error())
		return nil, err
	}

	// build and return the value
	return &models.BcBlock{
		Hash:      raw.Hash,
		Number:    models.Number(raw.Number),
		TimeStamp: graphql.Time{Time: time.Unix(int64(raw.Timestamp), 0)},
		TxHashes:  raw.Transactions,
	}, nil
}

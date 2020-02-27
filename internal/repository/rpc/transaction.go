package rpc

import (
	"fantomrocks-api/internal/models"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/graph-gophers/graphql-go"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

// Get a raw Transaction information for given tx hash.
func (rpc *Rpc) TransactionByHash(hash string) (*models.BcTransaction, error) {
	// unlock the source account
	rpc.log.Debugf("RPC->TransactionByHash(): Loading tx details for [%s]", hash)

	// container for raw data
	var raw struct {
		Hash      string        `json:"hash"`
		From      string        `json:"from"`
		To        *string       `json:"to"`
		Value     hexutil.Big   `json:"value"`
		Input     string        `json:"input"`
		Nonce     hexutil.Uint  `json:"nonce"`
		Gas       hexutil.Big   `json:"gas"`
		GasPrice  hexutil.Big   `json:"gasPrice"`
		BlockHash *string       `json:"blockHash"`
		TxIndex   *hexutil.Uint `json:"transactionIndex"`
	}

	// call for data
	err := rpc.Call(&raw, "ftm_getTransactionByHash", hash)
	if err != nil {
		rpc.log.Errorf("RPC->TransactionByHash(): Error! %s", err.Error())
		return nil, err
	}

	// get decoded base value
	value := big.Int(raw.Value)

	// get fee and gas related calculations
	var fee big.Int
	var gas big.Int
	gp := big.Int(raw.GasPrice)
	gl := big.Int(raw.Gas)

	// is there a block? get the receipt if we can
	if raw.BlockHash != nil {
		// get transaction receipt
		var rec struct {
			CumulativeGas hexutil.Big `json:"cumulativeGasUsed"`
			Gas           hexutil.Big `json:"gasUsed"`
		}

		// call for data
		err := rpc.Call(&rec, "eth_getTransactionReceipt", hash)
		if err != nil {
			rpc.log.Errorf("RPC->TransactionByHash(): Error! %s", err.Error())
			return nil, err
		}

		// calculate the fee
		gas = big.Int(rec.Gas)
		fee = *fee.Mul(&gp, &gas)
	}

	// index of the tx in the block (if any)
	var ix *int32
	if nil != raw.TxIndex {
		v := int32(*raw.TxIndex)
		ix = &v
	}

	// build and return the value
	return &models.BcTransaction{
		Hash:      raw.Hash,
		From:      raw.From,
		To:        raw.To,
		Value:     models.Amount{Decimal: decimal.NewFromBigInt(&value, 0)},
		Input:     raw.Input,
		Nonce:     uint(raw.Nonce),
		GasPrice:  models.Amount{Decimal: decimal.NewFromBigInt(&gp, 0)},
		GasLimit:  models.Amount{Decimal: decimal.NewFromBigInt(&gl, 0)},
		GasUsed:   models.Amount{Decimal: decimal.NewFromBigInt(&gas, 0)},
		Fee:       models.Amount{Decimal: decimal.NewFromBigInt(&fee, 0)},
		TxIndex:   ix,
		BlockHash: raw.BlockHash,
	}, nil
}

// Make a transfer of given amount of tokens from source account address to destination account address using given source account credentials.
func (rpc *Rpc) TransferTokens(fromAddr *models.Account, toAddr *models.Account, amount models.Amount) (*models.Transaction, error) {
	// unlock the source account
	rpc.log.Debugf("RPC->TransferTokens(): Sending %s tokens [%d => %d]", amount.ToHex(), fromAddr.Id, toAddr.Id)

	// prep transaction details
	tx := map[string]interface{}{
		"from":  fromAddr.Address,
		"to":    toAddr.Address,
		"value": amount.ToHex(),
	}

	// perform the call
	var txHash string
	err := rpc.Call(&txHash, "personal_sendTransaction", tx, fromAddr.Password)
	if err != nil {
		rpc.log.Errorf("RPC->TransferTokens(): Error! %s", err.Error())
		return nil, err
	}

	// unlock the source account
	rpc.log.Debugf("RPC->TransferTokens(): Tx [%d => %d] pending %s", fromAddr.Id, toAddr.Id, txHash)

	// return a valid
	return &models.Transaction{
		Id:          txHash,
		FromAccount: fromAddr,
		ToAccount:   toAddr,
		Amount:      &amount,
		TimeStamp:   &graphql.Time{Time: time.Now()},
	}, nil
}

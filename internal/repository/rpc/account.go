package rpc

import (
	"fantomrocks-api/internal/models"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

// Get Account balance from block-chain node
func (rpc *Rpc) AccountBalance(addr string) (*models.Amount, error) {
	// use RPC to make the call
	var balance string
	err := rpc.Call(&balance, "ftm_getBalance", addr, "latest")
	if err != nil {
		rpc.log.Errorf("RPC->AccountBalance(): Error [%s]", err.Error())
		return &models.Amount{}, err
	}

	// decode the response
	val, err := hexutil.DecodeBig(balance)
	if err != nil {
		rpc.log.Errorf("RPC->AccountBalance(): Can not get account balance for [%s]. %s", addr, err.Error())
		return &models.Amount{}, err
	}

	return &models.Amount{Decimal: decimal.NewFromBigInt(val, 0)}, nil
}

package rpc

import (
	"fantomrocks-api/internal/common"
	"fantomrocks-api/internal/models"
	"fantomrocks-api/internal/services"
	"github.com/ethereum/go-ethereum/rpc"
)

// BlockChain adapter interface definitions
type BlockChain interface {
	AccountBalance(string) (*models.Amount, error)
	TransactionByHash(string) (*models.BcTransaction, error)
	BlockByHash(string) (*models.BcBlock, error)
	TransferTokens(*models.Account, *models.Account, models.Amount) (*models.Transaction, error)
}

// Block-Chain RPC Adapter
type Rpc struct {
	log services.Logger
	*rpc.Client
}

// Prepare RPC client to be used to access block-chain node through it's com interface.
func NewRpc(cfg *common.Config, log services.Logger) (BlockChain, error) {
	// log actions
	log.Debugf("NewRpc(): Initializing RPC connection to Node [%s]", cfg.RpcUrl)

	// try to establish a connection
	client, err := rpc.Dial(cfg.RpcUrl)
	if err != nil {
		log.Criticalf("Can not connect to Node RPC end point. %s", err.Error())
		return nil, err
	}

	log.Debugf("NewRpc(): RPC adapter ready on [%s].", cfg.RpcUrl)
	return &Rpc{log: log, Client: client}, nil
}

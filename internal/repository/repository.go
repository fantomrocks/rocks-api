package repository

import (
	"fantomrocks-api/internal/common"
	dbRepo "fantomrocks-api/internal/repository/db"
	rpcRepo "fantomrocks-api/internal/repository/rpc"
	"fantomrocks-api/internal/services"
)

// Repository definition
type Repository struct {
	Db  dbRepo.DataStore
	Rpc rpcRepo.BlockChain
	Log services.Logger
}

// Create new Repository.
func NewRepository(cfg *common.Config, log services.Logger) (*Repository, error) {
	// get new data store adapter
	db, err := dbRepo.NewDB(cfg, log)
	if err != nil {
		log.Criticalf("NewEnv(): Database adapter not available, can not proceed! %s", err.Error())
		return nil, err
	}

	// get new block-chain node RPC adapter
	rpc, err := rpcRepo.NewRpc(cfg, log)
	if err != nil {
		log.Criticalf("NewEnv(): Database adapter not available, can not proceed! %s", err.Error())
		return nil, err
	}

	return &Repository{Db: db, Rpc: rpc, Log: log}, nil
}

package db

import (
	"fantomrocks-api/internal/common"
	"fantomrocks-api/internal/models"
	"fantomrocks-api/internal/services"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

// DataStore interface definition.
type DataStore interface {
	// accounts related
	AccountById(int) (*models.Account, error)
	AllAccounts() ([]*models.Account, error)
	RandomAccount() (*models.Account, error)
	RandomAccounts(count int, avoid []*models.Account) ([]*models.Account, error)

	// pairs related
	AllPairs() ([]*models.AccountPair, error)
	PairById(int) (*models.AccountPair, error)
	RandomPair() (*models.AccountPair, error)
}

// Database adapter
type DB struct {
	log services.Logger
	*sqlx.DB
}

// Get active adapter to a database holding additional data we need to serve the API.
func NewDB(cfg *common.Config, log services.Logger) (DataStore, error) {
	// log actions
	log.Debugf("NewDB(): Connecting database server [%s:%s/%s]", cfg.DbHost, cfg.DbPort, cfg.DbName)

	// try to open the connection
	db, err := sqlx.Open("pgx", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName))

	// do we have an error on obtaining the connection
	if err != nil {
		log.Criticalf("NewDB(): Fatal error, can not initialize DB layer! %s", err.Error())
		return nil, err
	}

	// try if the server is life
	if err = db.Ping(); err != nil {
		log.Criticalf("NewDB(): Fatal error, DB connection can not be established! %s", err.Error())
		return nil, err
	}

	// set additional connection details
	db.SetMaxOpenConns(cfg.DbMaxOpenConnections)

	// success
	log.Debugf("NewDB(): Database adapter ready.")
	return &DB{log: log, DB: db}, nil
}

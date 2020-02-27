package db

import (
	"fantomrocks-api/internal/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"time"
)

// define SQL queries used in service functions
const (
	sqlAccountById       string = "SELECT id, name, address, pwd FROM account WHERE id=$1"
	sqlAllAccounts       string = "SELECT id, name, address, pwd FROM account ORDER BY name"
	sqlAllAccountsExcept string = "SELECT id, name, address, pwd FROM account WHERE id NOT IN (?) ORDER BY name"
	sqlCountAccounts     string = `SELECT count(id) FROM account`
)

// Find account details by the account primary key.
func (db *DB) AccountById(id int) (*models.Account, error) {
	// make new Account
	acc := new(models.Account)

	// get the account
	err := db.Get(acc, sqlAccountById, id)
	if err != nil {
		db.log.Errorf("DB->AccountById(): Account #%d can not be loaded. %s", id, err.Error())
		return nil, err
	}

	// do we have the account?
	if 0 >= acc.Id {
		return nil, fmt.Errorf("DB->AccountById(): Account #%d not found", id)
	}
	return acc, err
}

// Get list of all accounts in the local database.
func (db *DB) AllAccounts() ([]*models.Account, error) {
	// make the container for results
	accounts := make([]*models.Account, 0)

	// try to get the data from database
	err := db.Select(&accounts, sqlAllAccounts)
	if err != nil {
		db.log.Errorf("DB->AllAccounts(): List of accounts can not be loaded. %s", err.Error())
	}

	return accounts, err
}

// Get single random account from the database; we don't expect gaps for deleted Accounts.
func (db *DB) RandomAccount() (*models.Account, error) {
	// get how many pairs we have
	var count int
	err := db.Get(&count, sqlCountAccounts)
	if err != nil {
		db.log.Errorf("DB->RandomAccount(): Can not count known Accounts. %s", err.Error())
		return nil, err
	}

	// inform
	db.log.Debugf("DB->RandomAccount(): Found %d Accounts to choose from.", count)

	// prep an empty Accounts
	acc := make([]*models.Account, 0)

	// do we have any pairs available?
	if 0 < count {
		// seed the random
		rand.Seed(time.Now().UnixNano())

		// we dont expect gaps so just pull the pair by random id
		err = db.Select(&acc, sqlAccountById, rand.Intn(count-1)+1)
		if err != nil {
			db.log.Errorf("DB->RandomAccount(): Can not get a random Account from database! %s", err.Error())
		}
	}

	// do we have any?
	if 0 == len(acc) {
		return nil, fmt.Errorf("random account not found")
	}

	// the pair could be send reversed by flip of a coin
	account := acc[0]
	return account, err
}

// Get list of <count> or less accounts skipping specified.
func (db *DB) RandomAccounts(count int, avoid []*models.Account) ([]*models.Account, error) {
	// prep accounts slice
	accounts := make([]*models.Account, 0)

	// make sure we return at least one element
	if 1 > count {
		db.log.Warningf("DB->RandomAccounts(): Expected to return at least one account. %d accounts requested!", count)
		count = 1
	}

	// limit the top
	if 50 < count {
		db.log.Warningf("DB->RandomAccounts(): Too many random accounts (%d) requested!", count)
		count = 50
	}

	// do we have any accounts to be avoided on the selection
	var err error
	if 0 < len(avoid) {
		// prep list of IDs of accounts to be avoided
		ids := make([]int, len(avoid))
		for i, acc := range avoid {
			ids[i] = int(acc.Id)
		}

		// expand binding to cover all the ids expected to be avoided
		query, args, err := sqlx.In(sqlAllAccountsExcept, ids)
		if err != nil {
			db.log.Errorf("DB->RandomAccounts(): Can not expand binding. %s", err.Error())
			return nil, err
		}

		// try top extract the data from database
		query = db.Rebind(query)
		db.log.Debugf("DB->RandomAccounts(): %s :: %v", query, args)
		err = db.Select(&accounts, query, args...)
	} else {
		// no limits; get them all
		err = db.Select(&accounts, sqlAllAccounts)
	}

	// did we succeeded with the query?
	if err != nil {
		db.log.Errorf("DB->RandomAccounts(): Accounts can not be loaded. %s", err.Error())
		return nil, err
	}

	// what we have
	db.log.Debugf("DB->RandomAccounts(): Found %d accounts.", len(accounts))

	// shuffle the slice we've got to extract random accounts
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(accounts), func(i, j int) { accounts[i], accounts[j] = accounts[j], accounts[i] })

	// extract subset of accounts with <count> entities if possible; we will return the full set otherwise
	if count < len(accounts) {
		accounts = accounts[:count]
	}

	// log & return
	db.log.Debugf("DB->RandomAccounts(): Requested %d, returning %d accounts.", count, len(accounts))
	return accounts, nil
}

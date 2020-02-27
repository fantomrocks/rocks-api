package db

import (
	"fantomrocks-api/internal/models"
	"math/rand"
	"time"
)

// define SQL queries used in service functions
const (
	sqlAllPairs string = `SELECT one.id as one_id, one.name as one_name, one.address as one_address, two.id as two_id, two.name as two_name, two.address as two_address 
								FROM account_pair JOIN account one ON one.id = account_pair.account_id_left JOIN account two ON two.id = account_pair.account_id_right 
    						ORDER BY one.name`
	sqlCountPairs      string = `SELECT count(id) FROM account_pair`
	sqlAccountPairById string = `SELECT one.id as one_id, one.name as one_name, one.address as one_address, two.id as two_id, two.name as two_name, two.address as two_address 
							FROM account_pair JOIN account one ON one.id = account_pair.account_id_left JOIN account two ON two.id = account_pair.account_id_right
							WHERE account_pair.id=$1`
)

// Get list of all account pairs from the database.
func (db *DB) AllPairs() ([]*models.AccountPair, error) {
	// make the container for results
	pairs := make([]*models.AccountPair, 0)

	// try to get the data from database
	rows, err := db.Query(sqlAllPairs)
	if err != nil {
		return pairs, err
	}

	// make sure the cursor is closed when we are done
	defer rows.Close()

	// loop rows
	for rows.Next() {
		// prep an empty Accounts
		one := new(models.Account)
		two := new(models.Account)

		// parse the query row and fill data elements
		err := rows.Scan(&one.Id, &one.Name, &one.Address, &two.Id, &two.Name, &two.Address)
		if err != nil {
			db.log.Errorf("DB->AllPairs(): Pairs row scan error! %s", err.Error())
		}

		// add new pair into the result set
		pairs = append(pairs, &models.AccountPair{One: one, Two: two})
	}

	err = rows.Err()
	return pairs, err
}

// Get list of all account pairs from the database.
func (db *DB) PairById(id int) (*models.AccountPair, error) {
	// inform
	db.log.Debugf("DB->PairById(): Loading Account Pair #%d.", id)

	// prep an empty Accounts
	one := new(models.Account)
	two := new(models.Account)

	// we dont expect gaps so just pull the pair by random id
	row := db.QueryRow(sqlAccountPairById, id)
	if err := row.Scan(&one.Id, &one.Name, &one.Address, &two.Id, &two.Name, &two.Address); err != nil {
		db.log.Errorf("DB->PairById(): Account Pair row scan error! %s", err.Error())
	}

	// the pair could be send reversed by flip of a coin
	return &models.AccountPair{One: one, Two: two}, nil
}

// Get random account pair from database; we don't expect gaps for deleted Pairs.
func (db *DB) RandomPair() (*models.AccountPair, error) {
	// get how many pairs we have
	var count int
	err := db.Get(&count, sqlCountPairs)
	if err != nil {
		db.log.Errorf("DB->RandomPair(): Can not count known pairs. %s", err.Error())
		return nil, err
	}

	// inform
	db.log.Debugf("DB->RandomPair(): Found %d Account Pairs to choose from.", count)

	// prep an empty Accounts
	one := new(models.Account)
	two := new(models.Account)

	// do we have any pairs available?
	if 0 < count {
		// seed the random
		rand.Seed(time.Now().UnixNano())

		// we dont expect gaps so just pull the pair by random id
		row := db.QueryRow(sqlAccountPairById, rand.Intn(count-1)+1)
		if err := row.Scan(&one.Id, &one.Name, &one.Address, &two.Id, &two.Name, &two.Address); err != nil {
			db.log.Errorf("DB->RandomPair(): Random Pair row scan error! %s", err.Error())
		}
	}

	// the pair could be send reversed by flip of a coin
	return &models.AccountPair{One: one, Two: two}, err
}

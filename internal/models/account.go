package models

// Define Account entity.
type Account struct {
	Id       int64
	Name     string
	Address  string
	Password string `db:"pwd"`
}

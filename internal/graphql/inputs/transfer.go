package inputs

import (
	"fantomrocks-api/internal/models"
	"github.com/graph-gophers/graphql-go"
)

// Defines structure to be used for describing an Account to Account test transfer.
type TransferInput struct {
	FromAccountId graphql.ID
	ToAccountId   graphql.ID
	Amount        models.Amount
}

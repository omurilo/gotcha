package entities

import (
	"errors"
	"time"
)

type TransactionType string

const (
	Debit  TransactionType = "d"
	Credit                 = "c"
)

type Transaction struct {
	Type        TransactionType `json:"tipo" bson:"type"`
	Value       uint64          `json:"valor" bson:"value"`
	Description *string         `json:"descricao" bson:"description"`
	CreatedAt   time.Time       `json:"realizada_em",omitempty bson:"createdat"`
	ClientId    uint64          `json:"-" bson:"client_id"`
}

type TransactionResponse struct {
	Limit   uint64 `json:"limite" bson:"limit"`
	Balance int64  `json:"saldo"  bson:"balance"`
}

func (t *Transaction) Validate() error {
	if t.Type == "" {
		return errors.New("The transaction type is required")
	}

	if t.Description == nil || len(*t.Description) < 1 || len(*t.Description) > 10 {
		return errors.New("Description must have between 1 and 10 characters")
	}

	if t.Type != Debit && t.Type != Credit {
		return errors.New("Type is not accepted")
	}

	return nil
}

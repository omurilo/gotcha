package entities

import (
	"errors"
	"math"
)

type Client struct {
	Id      uint64 `json:"id"     bson:"id"`
	Limit   uint64 `json:"limite" bson:"limit"`
	Balance int64  `json:"saldo"  bson:"balance"`
}

func (c *Client) Validate() error {
	if float64(c.Limit)-math.Abs(float64(c.Balance)) < 0 {
		return errors.New("Invalid transaction, user limit has been reached")
	}

	return nil
}

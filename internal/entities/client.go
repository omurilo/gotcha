package entities

import (
	"errors"
	"math"
)

type Client struct {
	Id      uint `json:"id"     bson:"id"`
	Limit   uint `json:"limite" bson:"limit"`
	Balance int  `json:"saldo"  bson:"balance"`
}

func (c *Client) Validate() error {
	if float64(c.Limit)-math.Abs(float64(c.Balance)) < 0 {
		return errors.New("Invalid transaction, user limit has been reached")
	}

	return nil
}

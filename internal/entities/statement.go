package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type balance struct {
	Total         int64     `bson:"total"          json:"total"`
	StatementDate time.Time `bson:"statement_date" json:"data_extrato"`
	Limit         uint64    `bson:"limit"          json:"limite"`
}

type Statement struct {
	Id           primitive.ObjectID `bson:"_id"          json:"-"`
	Balance      balance            `bson:"balance"      json:"saldo"`
	Transactions []Transaction      `bson:"transactions" json:"ultimas_transacoes"`
}

package repositories

import (
	"context"
	"errors"

	"github.com/omurilo/gotcha/internal/entities"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionsRepository struct {
	client *mongo.Client
}

func NewTransactionsRepository(client *mongo.Client) *TransactionsRepository {
	return &TransactionsRepository{client}
}

func (r *TransactionsRepository) Save(client *entities.Client, transaction *entities.Transaction) (*entities.TransactionResponse, error) {
	_, err := r.client.Database("rinha").Collection("transactions").InsertOne(context.TODO(), transaction)

	if err != nil {
		return nil, errors.New("Error has ocurred")
	}

	return &entities.TransactionResponse{client.Limit, client.Balance}, nil
}

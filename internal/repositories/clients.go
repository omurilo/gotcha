package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/omurilo/gotcha/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientsRepository struct {
	client *mongo.Client
}

func NewClientsRepository(client *mongo.Client) *ClientsRepository {
	return &ClientsRepository{client}
}

func (r *ClientsRepository) GetById(clientId uint64) (*entities.Client, error) {
	var client entities.Client
	err := r.client.Database("rinha").Collection("clients").FindOne(context.TODO(), bson.M{"id": clientId}).Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Client not found")
		}

		log.Fatal(err)
		return nil, err
	}

	return &client, nil
}

func (r *ClientsRepository) UpdateBalance(client *entities.Client) error {
	_, err := r.client.Database("rinha").
		Collection("clients").
		UpdateOne(context.TODO(), bson.M{"id": client.Id}, bson.M{"$set": bson.M{"balance": client.Balance}})

	return err
}

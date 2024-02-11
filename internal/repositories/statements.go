package repositories

import (
	"context"
	"log"

	"github.com/omurilo/gotcha/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatementsRepository struct {
	client *mongo.Client
}

func NewStatementsRepository(client *mongo.Client) *StatementsRepository {
	return &StatementsRepository{client}
}

func (r *StatementsRepository) Statement(client *entities.Client) (*entities.Statement, error) {
	pipeline := mongo.Pipeline{
		bson.D{{
			"$match", bson.D{{"id", client.Id}},
		}},
		bson.D{{
			"$addFields", bson.D{{"statement_date", "$$NOW"}},
		}},
		bson.D{{
			"$lookup", bson.D{
				{"from", "transactions"},
				{"localField", "id"},
				{"foreignField", "client_id"},
				{"as", "transactions"},
				{"pipeline", []bson.D{
					{{"$sort", bson.D{{"createdat", -1}}}},
					{{"$limit", 10}},
				}},
			},
		}},
		bson.D{{
			"$project", bson.D{
				{
					"balance", bson.D{
						{"total", "$balance"},
						{"statement_date", "$statement_date"},
						{"limit", "$limit"},
					},
				},
				{
					"transactions", "$transactions",
				},
			},
		}},
	}

	cursor, err := r.client.Database("rinha").Collection("clients").Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var results []entities.Statement
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &results[0], nil
}

package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/omurilo/gotcha/internal/entities"
)

type DbClient struct {
	Client *mongo.Client
}

func NewDbClient() *DbClient {
	DB_URL := os.Getenv("DB_URL")

	if DB_URL == "" {
		log.Fatal(
			"Could not be estabilish connection with database, please verify DB_URL environment variable",
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(DB_URL),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return &DbClient{client}
}

func (db *DbClient) InitDb() {
	client1 := entities.Client{Id: 1, Limit: 100000, Balance: 0}
	client2 := entities.Client{Id: 2, Limit: 80000, Balance: 0}
	client3 := entities.Client{Id: 3, Limit: 1000000, Balance: 0}
	client4 := entities.Client{Id: 4, Limit: 10000000, Balance: 0}
	client5 := entities.Client{Id: 5, Limit: 500000, Balance: 0}

	// Insert demo clients
	clientsCollection := db.Client.Database("rinha").Collection("clients")
	transactionsCollection := db.Client.Database("rinha").Collection("transactions")

	// clientsCollection.DeleteMany(context.TODO(), bson.D{})
	// transactionsCollection.DeleteMany(context.TODO(), bson.D{})

	models := []mongo.WriteModel{
		mongo.NewUpdateOneModel().
			SetFilter(bson.M{"id": 1}).
			SetUpdate(bson.D{{"$setOnInsert", client1}}).
			SetUpsert(true),
		mongo.NewUpdateOneModel().
			SetFilter(bson.D{{"id", 2}}).
			SetUpdate(bson.D{{"$setOnInsert", client2}}).
			SetUpsert(true),
		mongo.NewUpdateOneModel().
			SetFilter(bson.D{{"id", 3}}).
			SetUpdate(bson.D{{"$setOnInsert", client3}}).
			SetUpsert(true),
		mongo.NewUpdateOneModel().
			SetFilter(bson.D{{"id", 4}}).
			SetUpdate(bson.D{{"$setOnInsert", client4}}).
			SetUpsert(true),
		mongo.NewUpdateOneModel().
			SetFilter(bson.D{{"id", 5}}).
			SetUpdate(bson.D{{"$setOnInsert", client5}}).
			SetUpsert(true),
	}

  results, err := clientsCollection.BulkWrite(context.TODO(), models)
  if err != nil {
    fmt.Println("Error when insert accounts:", err)
  }

  fmt.Printf("Number of documents inserted: %d\n", results.UpsertedCount)

	clientsCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{"id", 1}}})
	transactionsCollection.Indexes().
		CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{"client_id", 1}, {"created_at", 1}}})
}

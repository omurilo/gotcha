package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id      uint64 `json:"id"     bson:"id"`
	Limit   uint64 `json:"limite" bson:"limit"`
	Balance int64  `json:"saldo"  bson:"balance"`
}

func DbConn() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://db:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func InitDb(client *mongo.Client) {
	user1 := User{Id: 1, Limit: 100000, Balance: 0}
	user2 := User{Id: 2, Limit: 80000, Balance: 0}
	user3 := User{Id: 3, Limit: 1000000, Balance: 0}
	user4 := User{Id: 4, Limit: 10000000, Balance: 0}
	user5 := User{Id: 5, Limit: 500000, Balance: 0}

	// Insert demo users
	collection := client.Database("rinha").Collection("users")
	models := []mongo.WriteModel{
		mongo.NewUpdateOneModel().
			SetFilter(bson.M{"id": 1}).
			SetUpdate(bson.D{{"$set", user1}}).
			SetUpsert(true),
		mongo.NewUpdateOneModel().
			SetFilter(bson.D{{"id", 2}}).
			SetUpdate(bson.D{{"$set", user2}}).
			SetUpsert(true),
		mongo.NewUpdateOneModel().
			SetFilter(bson.D{{"id", 3}}).
			SetUpdate(bson.D{{"$set", user3}}).
			SetUpsert(true),
		mongo.NewUpdateOneModel().
			SetFilter(bson.D{{"id", 4}}).
			SetUpdate(bson.D{{"$set", user4}}).
			SetUpsert(true),
		mongo.NewUpdateOneModel().
			SetFilter(bson.D{{"id", 5}}).
			SetUpdate(bson.D{{"$set", user5}}).
			SetUpsert(true),
	}
	collection.BulkWrite(
		context.TODO(),
		models,
	)

	fmt.Println("Demo users inserted")
}

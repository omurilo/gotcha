package router

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/omurilo/gotcha/cmd/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Transaction struct {
	Type        string    `json:"tipo" bson:"type"`
	Value       uint64    `json:"valor" bson:"value"`
	Description *string    `json:"descricao" bson:"description"`
	CreatedAt   time.Time `json:"realizada_em",omitempty bson:"createdat"`
	ClientId    uint64    `json:"-" bson:"client_id"`
}

type TransactionResponse struct {
	Limit   uint64 `json:"limite" bson:"limit"`
	Balance int64  `json:"saldo"  bson:"balance"`
}

var (
	usersCollection        *mongo.Collection
	transactionsCollection *mongo.Collection
)

func TransactionsRouter(client *mongo.Client, mux *http.ServeMux) {
	usersCollection = client.Database("rinha").Collection("users")
	transactionsCollection = client.Database("rinha").Collection("transactions")
	mux.HandleFunc("POST /clientes/{id}/transacoes", transactions)
}

func transactions(w http.ResponseWriter, r *http.Request) {
	clientId, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	var user database.User
	err = usersCollection.FindOne(context.TODO(), bson.M{"id": clientId}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User Not Found", http.StatusNotFound)
			return
		}

		log.Fatal(err)
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var body Transaction
	err = d.Decode(&body)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if body.Description == nil || len(*body.Description) < 1 || len(*body.Description) > 10 {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	if body.Type != "c" && body.Type != "d" {
		// fmt.Printf("Type is not accepted: %s\n", body.Type)
		http.Error(w, "Type is not accepted", http.StatusUnprocessableEntity)
		return
	}

	var newBalance int64
	if body.Type == "c" {
		newBalance = user.Balance + int64(body.Value)
	}

	if body.Type == "d" {
		newBalance = user.Balance - int64(body.Value)
	}

	if math.Abs(float64(newBalance)) > float64(user.Limit) {
		// fmt.Printf(
		// 	"User limit has been reached, limit: %s, balance: %s, value: %s and newBalance: %s\n",
		// 	user.Limit,
		// 	user.Balance,
		// 	body.Value,
		// 	newBalance,
		// )
		http.Error(w, "User limit has been reached", http.StatusUnprocessableEntity)
		return
	}

	usersCollection.UpdateOne(context.TODO(), bson.M{"id": clientId}, bson.D{{"$set", bson.M{"balance": newBalance}}})

	body.CreatedAt = time.Now()
	body.ClientId = clientId
	transactionsCollection.InsertOne(context.TODO(), body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TransactionResponse{user.Limit, newBalance})
}

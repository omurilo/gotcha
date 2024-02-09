package router

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/omurilo/gotcha/cmd/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func StatementsRouter(client *mongo.Client, mux *http.ServeMux) {
	usersCollection = client.Database("rinha").Collection("users")
	mux.HandleFunc("GET /clientes/{id}/extrato", statements)
}

func statements(w http.ResponseWriter, r *http.Request) {
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

	pipeline := mongo.Pipeline{
		bson.D{{
			"$match", bson.D{{"id", clientId}},
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
					bson.D{{"$sort", bson.D{{"createdat", -1}}}},
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

	var cursor *mongo.Cursor
	cursor, err = usersCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var results []Statement
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results[0])
}

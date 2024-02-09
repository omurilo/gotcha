package rinha

import (
	"context"
	"log"
	"net/http"

	"github.com/omurilo/gotcha/cmd/database"
	"github.com/omurilo/gotcha/cmd/router"
)

func Rinha() {
	client := database.DbConn()

	defer client.Disconnect(context.Background())

	// Insert demo users
	database.InitDb(client)

	mux := http.NewServeMux()
	router.TransactionsRouter(client, mux)
	router.StatementsRouter(client, mux)

	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}
}

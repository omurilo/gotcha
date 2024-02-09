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
	router.TransactionsRouter(client)
	router.StatementsRouter(client)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}

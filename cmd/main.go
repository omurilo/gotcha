package rinha

import (
	"context"
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

	http.ListenAndServe(":80", nil)
}

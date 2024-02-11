package main

import (
	"context"
	"log"
	"net/http"

	"github.com/omurilo/gotcha/infra/database"
	"github.com/omurilo/gotcha/internal/repositories"
	"github.com/omurilo/gotcha/internal/server"
	"github.com/omurilo/gotcha/internal/services"
)

func main() {
	dbClient := database.NewDbClient()

	defer dbClient.Client.Disconnect(context.Background())

	// Insert demo users
	dbClient.InitDb()

	clientsRepository := repositories.NewClientsRepository(dbClient.Client)
	transactionsRepository := repositories.NewTransactionsRepository(dbClient.Client)
	statementsRepository := repositories.NewStatementsRepository(dbClient.Client)

	transactionsService := services.NewTransactionsService(clientsRepository, transactionsRepository)
	statementsService := services.NewStatementsService(clientsRepository, statementsRepository)

	httpServer := server.NewHttpServer(statementsService, transactionsService)

	err := http.ListenAndServe(":80", httpServer.Instance)
	if err != nil {
		log.Fatal(err)
	}
}

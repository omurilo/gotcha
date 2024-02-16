package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/omurilo/gotcha/infra/database"
	"github.com/omurilo/gotcha/internal/repositories"
	"github.com/omurilo/gotcha/internal/server"
	"github.com/omurilo/gotcha/internal/services"
)

func main() {
	PORT, ok := os.LookupEnv("PORT")
	if !ok {
		PORT = "80"
	}

	dbClient := database.NewDbClient()

	defer dbClient.Client.Disconnect(context.Background())

	// Insert demo users
	dbClient.InitDb()

	clientsRepository := repositories.NewClientsRepository(dbClient.Client)
	transactionsRepository := repositories.NewTransactionsRepository(dbClient.Client)
	statementsRepository := repositories.NewStatementsRepository(dbClient.Client)

	transactionsService := services.NewTransactionsService(
		clientsRepository,
		transactionsRepository,
	)
	statementsService := services.NewStatementsService(clientsRepository, statementsRepository)

	httpServer := server.NewHttpServer(statementsService, transactionsService)

	panic(http.ListenAndServe(fmt.Sprintf(":%s", PORT), httpServer.Instance))
}

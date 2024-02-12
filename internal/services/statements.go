package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/omurilo/gotcha/internal/entities"
	"github.com/omurilo/gotcha/internal/repositories"
)

type StatementsService struct {
	clientsRepository    *repositories.ClientsRepository
	statementsRepository *repositories.StatementsRepository
}

func NewStatementsService(
	clientsRepository *repositories.ClientsRepository,
	statementsRepository *repositories.StatementsRepository,
) *StatementsService {
	return &StatementsService{clientsRepository, statementsRepository}
}

func (s *StatementsService) Statement(w http.ResponseWriter, r *http.Request) {
	clientId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Fatal(err)
		return
	}

	var client *entities.Client
	client, err = s.clientsRepository.GetById(clientId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var statement *entities.Statement
	statement, err = s.statementsRepository.Statement(client)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statement)
}

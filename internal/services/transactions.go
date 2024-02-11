package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/omurilo/gotcha/internal/entities"
	"github.com/omurilo/gotcha/internal/repositories"
)

type TransactionsService struct {
	clientsRepository      *repositories.ClientsRepository
	transactionsRepository *repositories.TransactionsRepository
}

func NewTransactionsService(
	clientsRepository *repositories.ClientsRepository,
	transactionsRepository *repositories.TransactionsRepository,
) *TransactionsService {
	return &TransactionsService{clientsRepository, transactionsRepository}
}

func (s *TransactionsService) Create(w http.ResponseWriter, r *http.Request) {
	clientId, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	var client *entities.Client
	client, err = s.clientsRepository.GetById(clientId)
	if err != nil {
		http.Error(w, "client Not Found", http.StatusNotFound)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var body entities.Transaction
	err = d.Decode(&body)
	if err != nil {
		// fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err = body.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if body.Type == "c" {
		client.Balance = client.Balance + int64(body.Value)
	}

	if body.Type == "d" {
		client.Balance = client.Balance - int64(body.Value)
	}

	if err = client.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = s.clientsRepository.UpdateBalance(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body.CreatedAt = time.Now()
	body.ClientId = client.Id

	var response *entities.TransactionResponse
	response, err = s.transactionsRepository.Save(client, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

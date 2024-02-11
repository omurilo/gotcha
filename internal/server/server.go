package server

import (
	"net/http"

	"github.com/omurilo/gotcha/internal/services"
)

type HttpServer struct {
	Instance *http.ServeMux
}

func NewHttpServer(s *services.StatementsService, t *services.TransactionsService) *HttpServer {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /clientes/{id}/extrato", s.Statement)
	mux.HandleFunc("POST /clientes/{id}/transacoes", t.Create)

	return &HttpServer{mux}
}

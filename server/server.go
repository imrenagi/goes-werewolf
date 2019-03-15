package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"context"
)

type Processor interface {
	Process(ctx context.Context, text string) (string, error)
}

func NewServer() *Server {
	s := &Server{
		Router: mux.NewRouter(),
	}
	s.routes()

	return s
}

type Server struct {
	Router    *mux.Router
	Processor Processor
}

func (s *Server) healthCheckHandle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("ok"))
		return
	}
}

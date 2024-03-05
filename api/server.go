package api

import (
	"fmt"
	"github.com/artemsmotritel/oktion/storage"
	"log"
	"net/http"
)

type Server struct {
	listenAddress string
	store         storage.Storage
}

func NewServer(listenAddress string, store storage.Storage) *Server {
	return &Server{
		listenAddress: listenAddress,
		store:         store,
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.listenAddress, s.newConfiguredRouter())
}

func (s *Server) newConfiguredRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", s.handleGetUsers)
	mux.HandleFunc("POST /users", s.handleCreateUser)
	mux.HandleFunc("GET /users/{id}", s.handleGetUserByID)
	mux.HandleFunc("PUT /users/{id}", s.handleUpdateUser)
	mux.HandleFunc("DELETE /users/{id}", s.handleDeleteUser)
	mux.HandleFunc("GET /auctions", s.handleGetAuctions)
	mux.HandleFunc("GET /auctions/{id}", s.handleGetAuctionByID)
	mux.HandleFunc("POST /auctions", s.handleCreateAuction)
	mux.HandleFunc("DELETE /auctions/{id}", s.handleDeleteAuction)

	return loggingMiddleware(mux)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := fmt.Sprintf("New Request: method - %s, url - %s", r.Method, r.URL.Path)
		log.Println(l)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (s *Server) badRequestError(w http.ResponseWriter, _ *http.Request, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

func (s *Server) internalError(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Something went very wrong at our part...", http.StatusInternalServerError)
}

func extractUserIDFromCookie(r *http.Request) (int64, error) {
	return 1, nil
}

package api

import (
	"encoding/json"
	"fmt"
	"github.com/artemsmotritel/oktion/storage"
	"github.com/artemsmotritel/oktion/types"
	"log"
	"net/http"
	"strconv"
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
	mux.HandleFunc("/foo", func(writer http.ResponseWriter, request *http.Request) {
		if err := json.NewEncoder(writer).Encode(nil); err != nil {
			log.Fatal(err.Error())
		}
		writer.Write([]byte("ok"))
	})

	return mux
}

func (s *Server) handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad user id in path: %s", r.PathValue("id")))
		return
	}

	user, err := s.store.GetUserByID(id)

	if err != nil {
		s.internalError(w, r)
		return
	}

	if user == nil {
		s.handleNotFound(w, r)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Fatal(err.Error())
	}
}

func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.store.GetUsers()

	if err != nil {
		s.internalError(w, r)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Fatal(err.Error())
	}
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	bodyReader := json.NewDecoder(r.Body)
	var userRequest types.UserCreateRequest

	if err := bodyReader.Decode(&userRequest); err != nil {
		s.badRequestError(w, r, "Bad request body")
		return
	}

	user := types.MapUserCreateRequest(userRequest)

	if err := s.store.SaveUser(user); err != nil {
		s.internalError(w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) badRequestError(w http.ResponseWriter, _ *http.Request, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "text/plain")
	if _, err := w.Write([]byte(message)); err != nil {
		log.Fatal(err.Error())
	}
}

func (s *Server) internalError(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Content-Type", "text/plain")
	if _, err := w.Write([]byte("Something went very wrong at our part...")); err != nil {
		log.Fatal(err.Error())
	}
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad user id in path: %s", r.PathValue("id")))
		return
	}

	bodyReader := json.NewDecoder(r.Body)
	var userRequest types.UserUpdateRequest

	if err := bodyReader.Decode(&userRequest); err != nil {
		s.badRequestError(w, r, "Bad request body")
		return
	}

	if err := s.store.UpdateUser(id, &userRequest); err != nil {
		s.internalError(w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad user id in path: %s", r.PathValue("id")))
		return
	}

	if err := s.store.DeleteUser(id); err != nil {
		s.internalError(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

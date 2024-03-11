package api

import (
	"encoding/json"
	"fmt"
	"github.com/artemsmotritel/oktion/types"
	"net/http"
	"strconv"
)

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
	if err = json.NewEncoder(w).Encode(user); err != nil {
		s.logger.Println("ERROR: ", err.Error())
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
	if err = json.NewEncoder(w).Encode(users); err != nil {
		s.logger.Println("ERROR: ", err.Error())
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

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad user id in path: %s", r.PathValue("id")))
		return
	}

	bodyReader := json.NewDecoder(r.Body)
	var userRequest types.UserUpdateRequest

	if err = bodyReader.Decode(&userRequest); err != nil {
		s.badRequestError(w, r, "Bad request body")
		return
	}

	if err = s.store.UpdateUser(id, &userRequest); err != nil {
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

	if err = s.store.DeleteUser(id); err != nil {
		s.internalError(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

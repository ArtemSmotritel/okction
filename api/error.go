package api

import (
	"github.com/artemsmotritel/oktion/templates"
	"net/http"
)

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	handler := templates.NewNotFoundPageHandler()
	handler.ServeHTTP(w, r)
}

func (s *Server) badRequestError(w http.ResponseWriter, _ *http.Request, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

func (s *Server) internalError(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Something went very wrong at our part...", http.StatusInternalServerError)
}

func (s *Server) handleUnauthorized(w http.ResponseWriter, r *http.Request) {
	handler := templates.NewUnauthorizedPageHandler()
	handler.ServeHTTP(w, r)
}

func (s *Server) handleForbidden(w http.ResponseWriter, r *http.Request) {
	handler := templates.NewForbiddenPageHandler()
	handler.ServeHTTP(w, r)
}

package api

import (
	"github.com/artemsmotritel/oktion/templates"
	"net/http"
)

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	handler := templates.NewErrorPageHandler(templates.NotFound)
	handler.ServeHTTP(w, r)
}

func (s *Server) badRequestError(w http.ResponseWriter, _ *http.Request, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

func (s *Server) internalError(w http.ResponseWriter, r *http.Request) {
	handler := templates.NewErrorPageHandler(templates.InternalServerError)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleUnauthorized(w http.ResponseWriter, r *http.Request) {
	handler := templates.NewErrorPageHandler(templates.Unauthorized)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleForbidden(w http.ResponseWriter, r *http.Request) {
	handler := templates.NewErrorPageHandler(templates.Forbidden)
	handler.ServeHTTP(w, r)
}

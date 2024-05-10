package api

import (
	"github.com/artemsmotritel/oktion/templates"
	"net/http"
)

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Retarget", "#main")
	w.Header().Set("HX-Reswap", "outerHTML")
	handler := templates.NewErrorPageHandler(templates.NotFound)
	w.WriteHeader(http.StatusNotFound)
	handler.ServeHTTP(w, r)
}

func (s *Server) badRequestError(w http.ResponseWriter, _ *http.Request, message string) {

	http.Error(w, message, http.StatusBadRequest)
}

func (s *Server) internalError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Retarget", "#main")
	w.Header().Set("HX-Reswap", "outerHTML")
	handler := templates.NewErrorPageHandler(templates.InternalServerError)
	w.WriteHeader(http.StatusInternalServerError)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleUnauthorized(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Retarget", "#main")
	w.Header().Set("HX-Reswap", "outerHTML")
	handler := templates.NewErrorPageHandler(templates.Unauthorized)
	w.WriteHeader(http.StatusUnauthorized)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleForbidden(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Retarget", "#main")
	w.Header().Set("HX-Reswap", "outerHTML")
	handler := templates.NewErrorPageHandler(templates.Forbidden)
	w.WriteHeader(http.StatusForbidden)
	handler.ServeHTTP(w, r)
}

func (s *Server) statusConflict(w http.ResponseWriter, r *http.Request, message string) {
	w.Header().Set("HX-Retarget", "#main")
	w.Header().Set("HX-Reswap", "outerHTML")
	handler := templates.NewErrorPageWithMessageHandler(templates.StatusConflict, message)
	w.WriteHeader(http.StatusConflict)
	handler.ServeHTTP(w, r)
}

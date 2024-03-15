package api

import (
	"github.com/artemsmotritel/oktion/templates"
	"github.com/artemsmotritel/oktion/utils"
	"net/http"
)

func (s *Server) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	hxBoosted, err := utils.ExtractValueFromContext[bool](r.Context(), "hxBoosted")
	if err != nil {
		hxBoosted = false
	}

	handler := templates.NewProfilePageHandler(!hxBoosted)
	handler.ServeHTTP(w, r)
}

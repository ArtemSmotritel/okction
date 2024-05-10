package api

import (
	"fmt"
	"net/http"
	"strconv"
)

func (s *Server) handleGetCategoryAuctions(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad category id in path: %s", r.PathValue("id")))
		return
	}

	w.Header().Set("Location", "/auctions?category="+strconv.FormatInt(id, 10))
	w.WriteHeader(http.StatusSeeOther)
}

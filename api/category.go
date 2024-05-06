package api

import (
	"errors"
	"fmt"
	"github.com/artemsmotritel/oktion/types"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
)

func (s *Server) handleGetCategoryAuctions(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad category id in path: %s", r.PathValue("id")))
		return
	}

	filterBuilder := types.NewAuctionFilterBuilder().WithCategoryId(id)

	_, err = s.store.GetAuctions(filterBuilder.Build())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// TODO return either a 404 or an empty list with a special message
			return
		}
		s.internalError(w, r)
		return
	}

	// TODO finish UI part
}

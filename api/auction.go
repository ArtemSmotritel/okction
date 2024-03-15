package api

import (
	"encoding/json"
	"fmt"
	"github.com/artemsmotritel/oktion/templates"
	"github.com/artemsmotritel/oktion/types"
	"net/http"
	"strconv"
)

func (s *Server) handleGetAuctions(w http.ResponseWriter, r *http.Request) {
	auctions, err := s.store.GetAuctions()

	if err != nil {
		// TODO introduce logging to error responses
		s.internalError(w, r)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(auctions); err != nil {
		s.logger.Println("ERROR: ", err.Error())
	}
}

func (s *Server) handleGetAuctionByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction id in path: %s", r.PathValue("id")))
		return
	}

	auction, err := s.store.GetAuctionByID(id)

	if err != nil {
		s.internalError(w, r)
		return
	}

	if auction == nil {
		s.handleNotFound(w, r)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(auction); err != nil {
		s.logger.Println("ERROR: ", err.Error())
	}
}

func (s *Server) handleCreateAuction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.badRequestError(w, r, "Couldn't parse form")
		return
	}

	ownerID, err := extractUserIDFromCookie(r)

	if err != nil {
		s.badRequestError(w, r, "Invalid cookie")
	}

	auction, err := types.MapAuctionCreateRequest(r.Form, ownerID)

	if err != nil {
		s.badRequestError(w, r, "Bad form request: "+err.Error())
		return
	}

	savedAuction, err := s.store.SaveAuction(auction)
	if err != nil {
		s.internalError(w, r)
		return
	}

	w.Header().Add("HX-Push-Url", fmt.Sprintf("/auctions/%d/edit", savedAuction.ID))
	renderer := templates.NewEditAuctionPageRenderer(savedAuction, []types.AuctionLot{
		{
			ID:   1,
			Name: "Lot 1",
		},
		{
			ID:   2,
			Name: "Lot 2",
		},
	})
	w.WriteHeader(http.StatusCreated)
	renderer.ServeHTTP(w, r)
}

func (s *Server) handleDeleteAuction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction id in path: %s", r.PathValue("id")))
		return
	}

	if err = s.store.DeleteAuction(id); err != nil {
		s.internalError(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

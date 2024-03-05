package api

import (
	"encoding/json"
	"fmt"
	"github.com/artemsmotritel/oktion/types"
	"log"
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
		log.Fatal(err.Error())
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
	if err := json.NewEncoder(w).Encode(auction); err != nil {
		log.Fatal(err.Error())
	}
}

func (s *Server) handleCreateAuction(w http.ResponseWriter, r *http.Request) {
	bodyReader := json.NewDecoder(r.Body)
	var auctionRequest types.AuctionCreateRequest

	if err := bodyReader.Decode(&auctionRequest); err != nil {
		s.badRequestError(w, r, "Bad request body")
		return
	}

	ownerID, err := extractUserIDFromCookie(r)

	if err != nil {
		s.badRequestError(w, r, "Invalid cookie")
	}

	auction := types.MapAuctionCreateRequest(auctionRequest, ownerID)

	if err = s.store.SaveAuction(auction); err != nil {
		s.internalError(w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
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

package api

import (
	"fmt"
	"github.com/artemsmotritel/oktion/templates"
	"github.com/artemsmotritel/oktion/types"
	"net/http"
	"strconv"
)

func (s *Server) handleCreateAuctionLot(w http.ResponseWriter, r *http.Request) {
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

	auctionLotCount, err := s.store.GetAuctionLotCount(auction.ID)
	auctionLotCount++
	if err != nil {
		s.internalError(w, r)
		return
	}

	savedAuctionLot, err := s.store.SaveAuctionLot(&types.AuctionLot{
		AuctionID: auction.ID,
		Name:      fmt.Sprintf("Lot %d", auctionLotCount),
	})
	if err != nil {
		s.internalError(w, r)
		return
	}

	handler := templates.NewAuctionLotListItemHandler(savedAuctionLot)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleEditAuctionLot(w http.ResponseWriter, r *http.Request) {
	_, err := strconv.ParseInt(r.PathValue("auctionId"), 10, 64)
	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction id in path: %s", r.PathValue("auctionId")))
		return
	}

	lotId, err := strconv.ParseInt(r.PathValue("lotId"), 10, 64)
	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction lot id in path: %s", r.PathValue("lotId")))
		return
	}

	auctionLot, err := s.store.GetAuctionLotByID(lotId)
	if err != nil {
		s.internalError(w, r)
		return
	}

	handler := templates.NewAuctionLotEditPageHandler(auctionLot)
	handler.ServeHTTP(w, r)
}

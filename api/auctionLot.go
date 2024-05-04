package api

import (
	"fmt"
	"github.com/artemsmotritel/oktion/templates"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/validation"
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

func (s *Server) handleUpdateAuctionLot(w http.ResponseWriter, r *http.Request) {
	auctionId, err := strconv.ParseInt(r.PathValue("auctionId"), 10, 64)
	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction id in path: %s", r.PathValue("auctionId")))
		return
	}

	lotId, err := strconv.ParseInt(r.PathValue("lotId"), 10, 64)
	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction lot id in path: %s", r.PathValue("lotId")))
		return
	}

	if err = r.ParseForm(); err != nil {
		s.badRequestError(w, r, err.Error())
		return
	}

	updateRequest, err := types.NewAuctionLotUpdateRequest(r.Form, lotId, auctionId)
	if err != nil {
		s.internalError(w, r)
		return
	}

	validator := validation.NewAuctionLotUpdateValidator(updateRequest)
	ok, err := validator.Validate()

	if err != nil {
		s.internalError(w, r)
		return
	}

	if !ok {
		auctionLotWithBadInfo := &types.AuctionLot{
			ID:           lotId,
			AuctionID:    auctionId,
			Name:         updateRequest.Name,
			Description:  updateRequest.Description,
			IsActive:     true,
			MinimalBid:   updateRequest.MinimalBid,
			ReservePrice: updateRequest.ReservePrice,
			BinPrice:     updateRequest.BinPrice,
		}
		// TODO: handle not 2xx status codes as intended
		//w.WriteHeader(http.StatusBadRequest)
		handler := templates.NewAuctionLotEditFormErrorBadRequestHandler(auctionLotWithBadInfo, validator.Errors)
		handler.ServeHTTP(w, r)
		return
	}

	auctionLot, err := s.store.UpdateAuctionLot(lotId, updateRequest)
	if err != nil {
		s.internalError(w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
	handler := templates.NewAuctionLotEditFormHandler(auctionLot)
	handler.ServeHTTP(w, r)
}

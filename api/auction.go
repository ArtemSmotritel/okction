package api

import (
	"encoding/json"
	"fmt"
	"github.com/artemsmotritel/oktion/templates"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
	"github.com/artemsmotritel/oktion/validation"
	"net/http"
	"strconv"
)

func (s *Server) handleNewAuction(w http.ResponseWriter, r *http.Request) {
	handler := templates.NewCreateAuctionPageHandler()
	handler.ServeHTTP(w, r)
}

func (s *Server) handleEditAuction(w http.ResponseWriter, r *http.Request) {
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

	auctionLots, err := s.store.GetAuctionLotsByAuctionID(auction.ID)
	if err != nil {
		s.internalError(w, r)
		return
	}

	handler := templates.NewEditAuctionPageHandler(auction, auctionLots)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleGetMyAuctions(w http.ResponseWriter, r *http.Request) {
	_, err := extractUserIDFromCookie(r)

	if err != nil {
		s.badRequestError(w, r, "not authorized")
		return
	}

	ownerId, err := utils.ExtractValueFromContext[int64](r.Context(), "userId")
	if err != nil {
		s.handleUnauthorized(w, r)
		return
	}

	auctions, err := s.store.GetAuctionsByOwnerId(ownerId)

	if err != nil {
		s.internalError(w, r)
		return
	}

	handler := templates.NewMyAuctionsPageHandler(auctions)
	handler.ServeHTTP(w, r)
}

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

	w.Header().Add("HX-Push-Url", fmt.Sprintf("/my-auctions/%d/edit", savedAuction.ID))
	handler := templates.NewEditAuctionPageHandler(savedAuction, []types.AuctionLot{})
	w.WriteHeader(http.StatusCreated)
	handler.ServeHTTP(w, r)
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

func (s *Server) handleUpdateAuction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction id in path: %s", r.PathValue("id")))
		return
	}

	if err = r.ParseForm(); err != nil {
		s.badRequestError(w, r, err.Error())
		return
	}

	updateRequest := types.NewAuctionUpdateRequest(r.Form, id)
	validator := validation.NewAuctionUpdateValidator(updateRequest)
	ok, err := validator.Validate()
	if err != nil {
		s.internalError(w, r)
		return
	}

	if !ok {
		// TODO handle bad request
		s.badRequestError(w, r, "Invalid update request")
		return
	}

	updatedAuction, err := s.store.UpdateAuction(updateRequest)
	if err != nil {
		s.internalError(w, r)
		return
	}

	auctionLots, err := s.store.GetAuctionLotsByAuctionID(updatedAuction.ID)
	if err != nil {
		s.internalError(w, r)
		return
	}

	handler := templates.NewEditAuctionPageHandler(updatedAuction, auctionLots)
	w.WriteHeader(http.StatusCreated)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleArchiveAuction(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ExtractValueFromContext[int64](r.Context(), "userId")
	if err != nil {
		// TODO : make user there is userId in each protected request handler
		s.badRequestError(w, r, "Not authorized")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction id in path: %s", r.PathValue("id")))
		return
	}

	if err = s.store.SetAuctionActiveStatus(id, false); err != nil {
		s.internalError(w, r)
		return
	}

	auctions, err := s.store.GetAuctionsByOwnerId(userId)
	if err != nil {
		s.internalError(w, r)
		return
	}

	w.Header().Set("HX-Retarget", "#main")
	w.Header().Set("HX-Reswap", "outerHTML")
	handler := templates.NewMyAuctionsPageHandler(auctions)
	handler.ServeHTTP(w, r)
}

func (s *Server) handleReinstateAuction(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ExtractValueFromContext[int64](r.Context(), "userId")
	if err != nil {
		// TODO : make user there is userId in each protected request handler
		s.badRequestError(w, r, "Not authorized")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, fmt.Sprintf("Bad auction id in path: %s", r.PathValue("id")))
		return
	}

	if err = s.store.SetAuctionActiveStatus(id, true); err != nil {
		s.internalError(w, r)
		return
	}

	auctions, err := s.store.GetAuctionsByOwnerId(userId)
	if err != nil {
		s.internalError(w, r)
		return
	}

	w.Header().Set("HX-Retarget", "#main")
	w.Header().Set("HX-Reswap", "outerHTML")
	handler := templates.NewMyAuctionsPageHandler(auctions)
	handler.ServeHTTP(w, r)
}

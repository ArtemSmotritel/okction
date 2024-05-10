package api

import (
	"fmt"
	"github.com/artemsmotritel/oktion/templates"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
	"github.com/artemsmotritel/oktion/validation"
	"net/http"
	"strconv"
)

func (s *Server) handleMakeBid(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ExtractValueFromContext[int64](r.Context(), "userId")
	if err != nil {
		// TODO : make user there is userId in each protected request handler
		s.badRequestError(w, r, "Not authorized")
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

	auctionLot, err := s.store.GetAuctionLotByID(lotId)
	if err != nil {
		s.internalError(w, r)
		return
	}

	request := types.NewBidMakeRequest(r.Form, lotId, userId)
	lotProvider := validation.SavedAuctionLotProvider{
		Lot: auctionLot,
	}
	validator := validation.NewBidMakeValidator(request, &lotProvider)
	ok, err := validator.Validate()
	if err != nil {
		s.internalError(w, r)
		return
	}

	// TODO move to validator
	canBeBidOn, err := s.store.CanBidOnAuctionLot(lotId)
	if err != nil {
		s.internalError(w, r)
		return
	}

	if !canBeBidOn {
		validator.Errors["value"] = "This lot can't be bid on"
	}

	if !ok || !canBeBidOn {
		handler := templates.NewMakeBidErrorBadRequestHandler(auctionLot, canBeBidOn, validator.Errors["value"])
		handler.ServeHTTP(w, r)
		return
	}

	bid, err := s.store.SaveAuctionLotBid(request)
	if err != nil {
		s.internalError(w, r)
		return
	}

	if request.Value.Cmp(auctionLot.BinPrice) >= 0 {
		if err = s.store.MarkBidAsWin(bid.ID); err != nil {
			s.internalError(w, r)
			return
		}

		if err = s.store.CloseAuctionLot(lotId); err != nil {
			s.internalError(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleGetMyBids(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ExtractValueFromContext[int64](r.Context(), "userId")
	if err != nil {
		// TODO : make user there is userId in each protected request handler
		s.badRequestError(w, r, "Not authorized")
		return
	}

	bids, err := s.store.GetUserBids(userId)
	if err != nil {
		s.internalError(w, r)
		return
	}

	handler := templates.NewUserBidsPageHandler(bids, r.Context())
	handler.ServeHTTP(w, r)
}

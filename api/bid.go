package api

import (
	"fmt"
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
	lotProvider := validation.SavedAuctionProvider{
		Lot: auctionLot,
	}
	validator := validation.NewBidMakeValidator(request, &lotProvider)
	ok, err := validator.Validate()
	if err != nil {
		s.internalError(w, r)
		return
	}

	if !ok {
		// TODO
		return
	}

	if request.Value.Cmp(auctionLot.BinPrice) >= 0 {
		// TODO finish by saving both bid and making the auction closed
	}

	bid, err := s.store.SaveAuctionLotBid(request)
	if err != nil {
		s.internalError(w, r)
		return
	}

	// TODO
	_ = bid.ID
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

	// TODO
	_ = bids[1]
}

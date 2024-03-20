package templates

import (
	"github.com/a-h/templ"
	"github.com/artemsmotritel/oktion/types"
	"net/http"
)

type AuctionLotListItemHandler struct {
	auctionLot *types.AuctionLot
}

func NewAuctionLotListItemHandler(auctionLot *types.AuctionLot) *AuctionLotListItemHandler {
	return &AuctionLotListItemHandler{
		auctionLot: auctionLot,
	}
}

func (a *AuctionLotListItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := templ.Handler(a.newAuctionLotListItem())
	handler.ServeHTTP(w, r)
}

func (a *AuctionLotListItemHandler) newAuctionLotListItem() templ.Component {
	return auctionLotListItem(a.auctionLot)
}

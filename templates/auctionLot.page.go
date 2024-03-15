package templates

import (
	"github.com/a-h/templ"
	"github.com/artemsmotritel/oktion/types"
	"net/http"
)

type AuctionLotListItemHandler struct {
	auctionId int64
}

func NewAuctionLotListItemHandler(auctionId int64) *AuctionLotListItemHandler {
	return &AuctionLotListItemHandler{
		auctionId: auctionId,
	}
}

func (a *AuctionLotListItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := templ.Handler(newAuctionLotListItem(a.auctionId))
	handler.ServeHTTP(w, r)
}

func newAuctionLotListItem(auctionId int64) templ.Component {
	return auctionLotListItem(&types.AuctionLot{
		Name: "Awesome Lot",
	}, auctionId)
}

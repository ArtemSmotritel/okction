package templates

import (
	"github.com/a-h/templ"
	"github.com/artemsmotritel/oktion/types"
	"net/http"
)

type AuctionLotListItemRenderer struct {
	auctionId int64
}

func NewAuctionLotListItemRenderer(auctionId int64) *AuctionLotListItemRenderer {
	return &AuctionLotListItemRenderer{
		auctionId: auctionId,
	}
}

func (a *AuctionLotListItemRenderer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := templ.Handler(newAuctionLotListItem(a.auctionId))
	handler.ServeHTTP(w, r)
}

func newAuctionLotListItem(auctionId int64) templ.Component {
	return auctionLotListItem(&types.AuctionLot{
		Name: "Awesome Lot",
	}, auctionId)
}

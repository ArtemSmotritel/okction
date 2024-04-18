package templates

import (
	"context"
	"github.com/a-h/templ"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
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

type AuctionLotEditPageHandler struct {
	auctionLot *types.AuctionLot
}

func NewAuctionLotEditPageHandler(auctionLot *types.AuctionLot) *AuctionLotEditPageHandler {
	return &AuctionLotEditPageHandler{
		auctionLot: auctionLot,
	}
}

func (a *AuctionLotEditPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := templ.Handler(a.newAuctionLotEditPage(r.Context()))
	handler.ServeHTTP(w, r)
}

func (a *AuctionLotEditPageHandler) newAuctionLotEditPage(ctx context.Context) templ.Component {
	hxBoosted, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")
	if err != nil {
		hxBoosted = false
	}

	if hxBoosted {
		return auctionLotEditPage(a.auctionLot)
	}

	isAuthorized, err := utils.ExtractValueFromContext[bool](ctx, "isAuthorized")
	if err != nil {
		isAuthorized = false
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(isAuthorized))
	builder.AppendComponent(auctionLotEditPage(a.auctionLot))
	builder.AppendComponent(mainFooter())

	return builder.Build()
}

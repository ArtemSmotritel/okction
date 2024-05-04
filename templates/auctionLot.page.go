package templates

import (
	"context"
	"github.com/a-h/templ"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
	"net/http"
)

func NewAuctionLotListItemHandler(auctionLot *types.AuctionLot) *utils.TemplateHandler {
	return &utils.TemplateHandler{
		Template: auctionLotListItem(auctionLot),
	}
}

type AuctionLotEditPageHandler struct {
	auctionLot *types.AuctionLot
}

func NewAuctionLotEditPageHandler(auctionLot *types.AuctionLot) *AuctionLotEditPageHandler {
	return &AuctionLotEditPageHandler{
		auctionLot: auctionLot,
	}
}

func NewAuctionLotEditFormHandler(auctionLot *types.AuctionLot) *utils.TemplateHandler {
	return &utils.TemplateHandler{
		Template: editAuctionLotForm(auctionLot, nil),
	}
}

func NewAuctionLotEditFormErrorBadRequestHandler(auctionLot *types.AuctionLot, errors map[string]string) *utils.TemplateHandler {
	if errors == nil {
		errors = make(map[string]string)
	}
	return &utils.TemplateHandler{
		Template: editAuctionLotForm(auctionLot, errors),
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

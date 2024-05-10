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
	categories []types.Category
}

func NewAuctionLotEditPageHandler(auctionLot *types.AuctionLot, categories []types.Category) *AuctionLotEditPageHandler {
	return &AuctionLotEditPageHandler{
		auctionLot: auctionLot,
		categories: categories,
	}
}

func NewAuctionLotEditFormHandler(auctionLot *types.AuctionLot, categories []types.Category) *utils.TemplateHandler {
	return &utils.TemplateHandler{
		Template: editAuctionLotForm(auctionLot, nil, categories),
	}
}

func NewAuctionLotEditFormErrorBadRequestHandler(auctionLot *types.AuctionLot, errors map[string]string, categories []types.Category) *utils.TemplateHandler {
	if errors == nil {
		errors = make(map[string]string)
	}
	return &utils.TemplateHandler{
		Template: editAuctionLotForm(auctionLot, errors, categories),
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
		return auctionLotEditPage(a.auctionLot, a.categories)
	}

	isAuthorized, err := utils.ExtractValueFromContext[bool](ctx, "isAuthorized")
	if err != nil {
		isAuthorized = false
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(isAuthorized))
	builder.AppendComponent(auctionLotEditPage(a.auctionLot, a.categories))
	builder.AppendComponent(mainFooter())

	return builder.Build()
}

func NewViewAuctionLotPageHandler(pageParam *types.AuctionLotViewPageParam, ctx context.Context) *utils.TemplateHandler {
	return &utils.TemplateHandler{
		Template: newViewAuctionLotPage(pageParam, ctx),
	}
}

func newViewAuctionLotPage(pageParam *types.AuctionLotViewPageParam, ctx context.Context) templ.Component {
	hxBoosted, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")
	if err != nil {
		hxBoosted = false
	}

	if hxBoosted {
		return viewAuctionLotPage(pageParam)
	}

	isAuthorized, err := utils.ExtractValueFromContext[bool](ctx, "isAuthorized")
	if err != nil {
		isAuthorized = false
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(isAuthorized))
	builder.AppendComponent(viewAuctionLotPage(pageParam))
	builder.AppendComponent(mainFooter())

	return builder.Build()
}

func NewSavedUnsavedAuctionLotButtonHandler(auctionId, lotId int64, doesUserFollowLot bool) *utils.TemplateHandler {
	return &utils.TemplateHandler{
		Template: followAuctionLotButton(auctionId, lotId, doesUserFollowLot),
	}
}

func NewMakeBidErrorBadRequestHandler(lot *types.AuctionLot, canBeBidOn bool, err string) *utils.TemplateHandler {
	return &utils.TemplateHandler{
		Template: makeBidOnAuctionLotForm(lot, canBeBidOn, err),
	}
}

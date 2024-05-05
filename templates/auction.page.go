package templates

import (
	"context"
	"github.com/a-h/templ"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
	"net/http"
)

type CreateAuctionPageHandler struct {
}

type EditAuctionPageHandler struct {
	auctionLots []types.AuctionLot
	auction     *types.Auction
}

type MyAuctionsPageHandler struct {
	auctions []types.Auction
}

func NewCreateAuctionPageHandler() *CreateAuctionPageHandler {
	return &CreateAuctionPageHandler{}
}

func NewEditAuctionPageHandler(auction *types.Auction, auctionLots []types.AuctionLot) *EditAuctionPageHandler {
	return &EditAuctionPageHandler{
		auctionLots: auctionLots,
		auction:     auction,
	}
}

func NewMyAuctionsPageHandler(auctions []types.Auction) *MyAuctionsPageHandler {
	return &MyAuctionsPageHandler{
		auctions: auctions,
	}
}

func (r *CreateAuctionPageHandler) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newCreateAuctionPage(re.Context()))
	handler.ServeHTTP(w, re)
}

func (r *EditAuctionPageHandler) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newEditAuctionPage(re.Context(), r))
	handler.ServeHTTP(w, re)
}

func (r *MyAuctionsPageHandler) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newMyAuctionsPage(re.Context(), r))
	handler.ServeHTTP(w, re)
}

func newEditAuctionPage(ctx context.Context, handler *EditAuctionPageHandler) templ.Component {
	hxBoosted, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")

	if err != nil {
		hxBoosted = false
	}

	if hxBoosted {
		return editAuctionPage(handler.auctionLots, handler.auction, nil)
	}

	isAuthorized, err := utils.ExtractValueFromContext[bool](ctx, "isAuthorized")

	if err != nil {
		isAuthorized = false
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(isAuthorized))
	builder.AppendComponent(editAuctionPage(handler.auctionLots, handler.auction, nil))
	builder.AppendComponent(mainFooter())

	return builder.Build()
}

func newCreateAuctionPage(ctx context.Context) templ.Component {
	hxBoosted, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")

	if err != nil {
		hxBoosted = false
	}

	if hxBoosted {
		return createAuctionPage()
	}

	isAuthorized, err := utils.ExtractValueFromContext[bool](ctx, "isAuthorized")

	if err != nil {
		isAuthorized = false
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(isAuthorized))
	builder.AppendComponent(createAuctionPage())
	builder.AppendComponent(mainFooter())

	return builder.Build()
}

func newMyAuctionsPage(ctx context.Context, handler *MyAuctionsPageHandler) templ.Component {
	hxBoosted, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")
	if err != nil {
		hxBoosted = false
	}

	if hxBoosted {
		return myAuctionsPage(handler.auctions)
	}

	isAuthorized, err := utils.ExtractValueFromContext[bool](ctx, "isAuthorized")
	if err != nil {
		isAuthorized = false
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(isAuthorized))
	builder.AppendComponent(myAuctionsPage(handler.auctions))
	builder.AppendComponent(mainFooter())

	return builder.Build()
}

func NewAuctionLotsListHandler(auctionLots []types.AuctionLot, auction *types.Auction) *utils.TemplateHandler {
	return &utils.TemplateHandler{
		Template: auctionLotsList(auction, auctionLots),
	}
}

func NewAuctionEditFormErrorBadRequestHandler(auction *types.Auction, errors map[string]string) *utils.TemplateHandler {
	return &utils.TemplateHandler{
		Template: createAuctionForm(false, auction, errors),
	}
}

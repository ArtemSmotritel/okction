package templates

import (
	"context"
	"github.com/a-h/templ"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
	"net/http"
)

type CreateAuctionPageRenderer struct {
}

type EditAuctionPageRenderer struct {
	auctionLots []types.AuctionLot
	auction     *types.Auction
}

type MyAuctionsPageRenderer struct {
	auctions []types.Auction
}

func NewCreateAuctionPageRenderer() *CreateAuctionPageRenderer {
	return &CreateAuctionPageRenderer{}
}

func NewEditAuctionPageRenderer(auction *types.Auction, auctionLots []types.AuctionLot) *EditAuctionPageRenderer {
	return &EditAuctionPageRenderer{
		auctionLots: auctionLots,
		auction:     auction,
	}
}

func NewMyAuctionsPageRenderer(auctions []types.Auction) *MyAuctionsPageRenderer {
	return &MyAuctionsPageRenderer{
		auctions: auctions,
	}
}

func (r *CreateAuctionPageRenderer) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newCreateAuctionPage(re.Context()))
	handler.ServeHTTP(w, re)
}

func (r *EditAuctionPageRenderer) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newEditAuctionPage(re.Context(), r))
	handler.ServeHTTP(w, re)
}

func (r *MyAuctionsPageRenderer) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newMyAuctionsPage(re.Context(), r))
	handler.ServeHTTP(w, re)
}

func newEditAuctionPage(ctx context.Context, renderer *EditAuctionPageRenderer) templ.Component {
	hxBoosted, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")

	if err != nil {
		hxBoosted = false
	}

	if hxBoosted {
		return editAuctionPage(renderer.auctionLots, renderer.auction)
	}

	isAuthorized, err := utils.ExtractValueFromContext[bool](ctx, "isAuthorized")

	if err != nil {
		isAuthorized = false
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(isAuthorized))
	builder.AppendComponent(editAuctionPage(renderer.auctionLots, renderer.auction))
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

func newMyAuctionsPage(ctx context.Context, renderer *MyAuctionsPageRenderer) templ.Component {
	hxBoosted, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")
	if err != nil {
		hxBoosted = false
	}

	if hxBoosted {
		return myAuctionsPage(renderer.auctions)
	}

	isAuthorized, err := utils.ExtractValueFromContext[bool](ctx, "isAuthorized")
	if err != nil {
		isAuthorized = false
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(isAuthorized))
	builder.AppendComponent(myAuctionsPage(renderer.auctions))
	builder.AppendComponent(mainFooter())

	return builder.Build()
}

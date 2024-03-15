package templates

import (
	"context"
	"github.com/a-h/templ"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
	"log"
	"net/http"
)

type CreateAuctionPageRenderer struct {
}

type EditAuctionPageRenderer struct {
	auctionLots []types.AuctionLot
	auction     *types.Auction
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

func (r *CreateAuctionPageRenderer) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newCreateAuctionPage(re.Context()))
	handler.ServeHTTP(w, re)
}

func (r *EditAuctionPageRenderer) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	log.Printf("header in renderer: %v \n", w.Header())
	handler := templ.Handler(newEditAuctionPage(re.Context(), r))
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
		return createAuction()
	}

	isAuthorized, err := utils.ExtractValueFromContext[bool](ctx, "isAuthorized")

	if err != nil {
		isAuthorized = false
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(isAuthorized))
	builder.AppendComponent(createAuction())
	builder.AppendComponent(mainFooter())

	return builder.Build()
}

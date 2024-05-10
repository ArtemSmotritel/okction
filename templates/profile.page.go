package templates

import (
	"context"
	"github.com/a-h/templ"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
	"net/http"
)

type ProfileMenuItem struct {
	Name string
	Link string
}

type ProfilePageHandler struct {
	menuItems []ProfileMenuItem
	user      *types.User
}

func NewProfilePageHandler(user *types.User) *ProfilePageHandler {
	return &ProfilePageHandler{
		menuItems: []ProfileMenuItem{
			{
				Name: "Your auctions",
				Link: "/my-auctions",
			},
			{
				Name: "Your Favorite lots",
				Link: "/my-saved-auction-lots",
			},
			{
				Name: "Your bids",
				Link: "/my-bids",
			},
		},
		user: user,
	}
}

func (h *ProfilePageHandler) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(h.newProfilePage(re.Context()))
	handler.ServeHTTP(w, re)
}

func (h *ProfilePageHandler) newProfilePage(ctx context.Context) templ.Component {
	isHTMXRequest, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")
	if err != nil {
		isHTMXRequest = false
	}

	if isHTMXRequest {
		return profile(h.menuItems, h.user)
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(ctx.Value("isAuthorized").(bool)))
	builder.AppendComponent(profile(h.menuItems, h.user))
	builder.AppendComponent(mainFooter())

	return builder.Build()
}

func NewEditProfileErrorBadRequestHandler(user *types.User, errors map[string]string) *utils.TemplateHandler {
	if errors == nil {
		errors = make(map[string]string)
	}
	return &utils.TemplateHandler{
		Template: profileEditForm(user, errors),
	}
}

func NewProfileFormHandler(user *types.User) *utils.TemplateHandler {
	return &utils.TemplateHandler{
		Template: profileEditForm(user, nil),
	}
}

func NewSavedAuctionLotsPageHandler(lots []types.AuctionLot, ctx context.Context) *utils.TemplateHandler {
	return &utils.TemplateHandler{
		Template: newSavedAuctionLotsPage(lots, ctx),
	}
}

func newSavedAuctionLotsPage(lots []types.AuctionLot, ctx context.Context) templ.Component {
	isHTMXRequest, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")
	if err != nil {
		isHTMXRequest = false
	}

	if isHTMXRequest {
		return savedAuctionLotsPage(lots)
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(ctx.Value("isAuthorized").(bool)))
	builder.AppendComponent(savedAuctionLotsPage(lots))
	builder.AppendComponent(mainFooter())

	return builder.Build()
}

func NewUserBidsPageHandler(bids []types.UserBid, ctx context.Context) *utils.TemplateHandler {
	return &utils.TemplateHandler{
		Template: newUserBidsPage(bids, ctx),
	}
}

func newUserBidsPage(bids []types.UserBid, ctx context.Context) templ.Component {
	isHTMXRequest, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")
	if err != nil {
		isHTMXRequest = false
	}

	if isHTMXRequest {
		return userBidsPage(bids)
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(ctx.Value("isAuthorized").(bool)))
	builder.AppendComponent(userBidsPage(bids))
	builder.AppendComponent(mainFooter())

	return builder.Build()
}

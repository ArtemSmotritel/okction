package templates

import (
	"context"
	"github.com/a-h/templ"
	"net/http"
)

type ProfileMenuItem struct {
	Name string
	Link string
}

type ProfilePageHandler struct {
	shouldBuildWholePage bool
}

func NewProfilePageHandler(shouldBuildWholePage bool) *ProfilePageHandler {
	return &ProfilePageHandler{
		shouldBuildWholePage: shouldBuildWholePage,
	}
}

func (r *ProfilePageHandler) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newProfilePage(re.Context(), r.shouldBuildWholePage))
	handler.ServeHTTP(w, re)
}

func newProfilePage(ctx context.Context, shouldBuildWholePage bool) templ.Component {
	items := []ProfileMenuItem{
		{
			Name: "Your auctions",
			Link: "/my-auctions",
		},
		{
			Name: "Your Favorite lots",
			Link: "/my-favorite-lots",
		},
		{
			Name: "Your bids",
			Link: "/my-bids",
		},
	}

	if shouldBuildWholePage {
		builder := NewHTMLPageBuilder(root)
		builder.AppendComponent(mainHeader(ctx.Value("isAuthorized").(bool)))
		builder.AppendComponent(profile(items))
		builder.AppendComponent(mainFooter())

		return builder.Build()
	}

	return profile(items)
}

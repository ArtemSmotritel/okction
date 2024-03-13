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

type ProfilePageRenderer struct {
	shouldBuildWholePage bool
}

func NewProfilePageRenderer(shouldBuildWholePage bool) *ProfilePageRenderer {
	return &ProfilePageRenderer{
		shouldBuildWholePage: shouldBuildWholePage,
	}
}

func (r *ProfilePageRenderer) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newProfilePage(re.Context(), r.shouldBuildWholePage))
	handler.ServeHTTP(w, re)
}

func newProfilePage(ctx context.Context, shouldBuildWholePage bool) templ.Component {
	items := []ProfileMenuItem{
		{
			Name: "My auctions",
			Link: "/my-lots",
		},
		{
			Name: "My Favorite lots",
			Link: "/my-favorite-lots",
		},
		{
			Name: "My bids",
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

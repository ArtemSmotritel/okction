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
}

func NewProfilePageRenderer() *ProfilePageRenderer {
	return &ProfilePageRenderer{}
}

func (r *ProfilePageRenderer) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newProfilePage(re.Context()))
	handler.ServeHTTP(w, re)
}

func newProfilePage(ctx context.Context) templ.Component {
	items := []ProfileMenuItem{{
		Name: "Favorite lots",
		Link: "/favorite-lots",
	}, {
		Name: "My auctions",
		Link: "/my-lots",
	},
	}

	builder := NewHTMLPageBuilder(Root)
	builder.AppendComponent(MainHeader(ctx.Value("isAuthorized").(bool)))
	builder.AppendComponent(profile(items))
	builder.AppendComponent(MainFooter())

	return builder.Build()
}

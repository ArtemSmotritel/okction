package templates

import (
	"context"
	"github.com/a-h/templ"
	"github.com/artemsmotritel/oktion/types"
	"net/http"
)

type IndexPageRenderer struct {
	categories []types.Category
}

func NewIndexPageRenderer(categories []types.Category) *IndexPageRenderer {
	return &IndexPageRenderer{
		categories: categories,
	}
}

func (r *IndexPageRenderer) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newIndexPage(r.categories, re.Context()))
	handler.ServeHTTP(w, re)
}

func newIndexPage(categories []types.Category, ctx context.Context) templ.Component {
	builder := NewHTMLPageBuilder(Root)
	builder.AppendComponent(MainHeader(ctx.Value("isAuthorized").(bool)))
	builder.AppendComponent(MainMain(categories...))
	builder.AppendComponent(MainFooter())

	return builder.Build()
}

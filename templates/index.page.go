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
	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(ctx.Value("isAuthorized").(bool)))
	builder.AppendComponent(mainMain(categories...))
	builder.AppendComponent(mainFooter())

	return builder.Build()
}

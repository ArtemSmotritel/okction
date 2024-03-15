package templates

import (
	"context"
	"github.com/a-h/templ"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
	"net/http"
)

type IndexPageHandler struct {
	categories []types.Category
}

func NewIndexPageHandler(categories []types.Category) *IndexPageHandler {
	return &IndexPageHandler{
		categories: categories,
	}
}

func (r *IndexPageHandler) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newIndexPage(r.categories, re.Context()))
	handler.ServeHTTP(w, re)
}

func newIndexPage(categories []types.Category, ctx context.Context) templ.Component {
	hxBoosted, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")

	if err != nil {
		hxBoosted = false
	}

	if hxBoosted {
		return mainMain(categories...)
	}

	isAuthorized, err := utils.ExtractValueFromContext[bool](ctx, "isAuthorized")

	if err != nil {
		isAuthorized = false
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(isAuthorized))
	builder.AppendComponent(mainMain(categories...))
	builder.AppendComponent(mainFooter())

	return builder.Build()
}

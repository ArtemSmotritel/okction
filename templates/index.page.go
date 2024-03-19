package templates

import (
	"context"
	"github.com/a-h/templ"
	"github.com/artemsmotritel/oktion/types"
	"github.com/artemsmotritel/oktion/utils"
	"net/http"
)

type IndexPageHandler struct {
	categories            []types.Category
	shouldRenderWholeBody bool
}

func NewIndexPageHandler(categories []types.Category) *IndexPageHandler {
	return &IndexPageHandler{
		categories:            categories,
		shouldRenderWholeBody: false,
	}
}

func NewIndexBodyHandler(categories []types.Category) *IndexPageHandler {
	return &IndexPageHandler{
		categories:            categories,
		shouldRenderWholeBody: true,
	}
}

func (h *IndexPageHandler) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(h.newIndexPage(h.categories, re.Context()))
	handler.ServeHTTP(w, re)
}

func (h *IndexPageHandler) newIndexPage(categories []types.Category, ctx context.Context) templ.Component {
	isAuthorized, err := utils.ExtractValueFromContext[bool](ctx, "isAuthorized")
	if err != nil {
		isAuthorized = false
	}

	var builder *HTMLPageBuilder

	if h.shouldRenderWholeBody {
		builder = NewHTMLPageBuilder(body)
	} else {
		builder = NewHTMLPageBuilder(root)
	}

	hxBoosted, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")
	if err != nil {
		hxBoosted = false
	}

	if hxBoosted && !h.shouldRenderWholeBody {
		return mainMain(categories...)
	}

	builder.AppendComponent(mainHeader(isAuthorized))
	builder.AppendComponent(mainMain(categories...))
	builder.AppendComponent(mainFooter())

	return builder.Build()
}

package templates

import (
	"context"
	"github.com/a-h/templ"
	"net/http"
)

type NotFoundPageRenderer struct {
}

func NewNotFoundPageRenderer() *NotFoundPageRenderer {
	return &NotFoundPageRenderer{}
}

func (r *NotFoundPageRenderer) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newNotFoundPage(re.Context()))
	handler.ServeHTTP(w, re)
}

func newNotFoundPage(ctx context.Context) templ.Component {
	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(ctx.Value("isAuthorized").(bool)))
	builder.AppendComponent(notFound())

	return builder.Build()
}

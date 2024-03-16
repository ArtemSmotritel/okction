package templates

import (
	"context"
	"github.com/a-h/templ"
	"net/http"
)

type NotFoundPageHandler struct {
}

func NewNotFoundPageHandler() *NotFoundPageHandler {
	return &NotFoundPageHandler{}
}

func (r *NotFoundPageHandler) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newNotFoundPage(re.Context()))
	handler.ServeHTTP(w, re)
}

func newNotFoundPage(ctx context.Context) templ.Component {
	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(ctx.Value("isAuthorized").(bool)))
	builder.AppendComponent(notFound())

	return builder.Build()
}

type UnauthorizedPageHandler struct {
}

func NewUnauthorizedPageHandler() *UnauthorizedPageHandler {
	return &UnauthorizedPageHandler{}
}

func (r *UnauthorizedPageHandler) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newUnauthorizedPage(re.Context()))
	handler.ServeHTTP(w, re)
}

func newUnauthorizedPage(ctx context.Context) templ.Component {
	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(ctx.Value("isAuthorized").(bool)))
	builder.AppendComponent(unauthorized())

	return builder.Build()
}

type ForbiddenPageHandler struct {
}

func NewForbiddenPageHandler() *ForbiddenPageHandler {
	return &ForbiddenPageHandler{}
}

func (r *ForbiddenPageHandler) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newForbiddenPage(re.Context()))
	handler.ServeHTTP(w, re)
}

func newForbiddenPage(ctx context.Context) templ.Component {
	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(mainHeader(ctx.Value("isAuthorized").(bool)))
	builder.AppendComponent(forbidden())

	return builder.Build()
}

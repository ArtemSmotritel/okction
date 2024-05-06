package templates

import (
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/artemsmotritel/oktion/utils"
	"net/http"
)

type ErrorCode int

const (
	NotFound            ErrorCode = 1
	Forbidden           ErrorCode = 2
	Unauthorized        ErrorCode = 3
	InternalServerError ErrorCode = 4
	StatusConflict      ErrorCode = 5
)

type ErrorPageHandler struct {
	template templ.Component
}

func NewErrorPageWithMessageHandler(errorCode ErrorCode, message string) *ErrorPageHandler {
	// TODO: maybe refactor to use Builder pattern
	var template templ.Component

	switch errorCode {
	case NotFound:
		template = notFound()
	case Forbidden:
		template = forbidden()
	case Unauthorized:
		template = unauthorized()
	case InternalServerError:
		template = internal()
	case StatusConflict:
		template = statusConflict(message)
	default:
		panic(fmt.Sprintf("unsupported error code was provided: %d", errorCode))
	}

	return &ErrorPageHandler{
		template: template,
	}
}

func NewErrorPageHandler(errorCode ErrorCode) *ErrorPageHandler {
	return NewErrorPageWithMessageHandler(errorCode, "")
}

func (r *ErrorPageHandler) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	hxBoosted, err := utils.ExtractValueFromContext[bool](re.Context(), "hxBoosted")
	if err != nil {
		hxBoosted = false
	}

	if hxBoosted {
		w.Header().Set("HX-Retarget", "body")
		w.Header().Set("HX-Reswap", "innerHTML")
	}

	handler := templ.Handler(newErrorPage(re.Context(), r.template))
	handler.ServeHTTP(w, re)
}

func newErrorPage(ctx context.Context, errorTemplate templ.Component) templ.Component {
	hxBoosted, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")
	if err != nil {
		hxBoosted = false
	}

	var builder *HTMLPageBuilder

	if hxBoosted {
		builder = NewHTMLPageBuilder(body)
	} else {
		builder = NewHTMLPageBuilder(root)
	}

	builder.AppendComponent(mainHeader(ctx.Value("isAuthorized").(bool)))
	builder.AppendComponent(errorTemplate)

	return builder.Build()
}

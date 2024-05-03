package templates

import (
	"context"
	"github.com/a-h/templ"
	"github.com/artemsmotritel/oktion/utils"
	"net/http"
)

type LoginPageHandler struct {
}

func NewLoginPageHandler() *LoginPageHandler {
	return &LoginPageHandler{}
}

func (r *LoginPageHandler) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newLoginPage(re.Context()))
	handler.ServeHTTP(w, re)
}

func newLoginPage(ctx context.Context) templ.Component {
	isHTMXRequest, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")
	if err != nil {
		isHTMXRequest = false
	}

	if isHTMXRequest {
		return login()
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(login())

	return builder.Build()
}

type SignUpPageHandler struct {
}

func NewSignUpPageHandler() *SignUpPageHandler {
	return &SignUpPageHandler{}
}

func (r *SignUpPageHandler) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handler := templ.Handler(newSignUpPage(re.Context()))
	handler.ServeHTTP(w, re)
}

func newSignUpPage(ctx context.Context) templ.Component {
	isHTMXRequest, err := utils.ExtractValueFromContext[bool](ctx, "hxBoosted")
	if err != nil {
		isHTMXRequest = false
	}

	if isHTMXRequest {
		return signUp()
	}

	builder := NewHTMLPageBuilder(root)
	builder.AppendComponent(signUp())

	return builder.Build()
}

func NewSignUpErrorBadRequestHandler(values map[string]string, errors map[string]string) *utils.TemplateHandler {
	if values == nil {
		values = make(map[string]string)
	}
	if errors == nil {
		errors = make(map[string]string)
	}
	return &utils.TemplateHandler{
		Template: signUpForm(values, errors),
	}
}

func NewLoginErrorBadRequestHandler(values map[string]string, errors map[string]string) *utils.TemplateHandler {
	if values == nil {
		values = make(map[string]string)
	}
	if errors == nil {
		errors = make(map[string]string)
	}
	return &utils.TemplateHandler{
		Template: loginForm(values, errors),
	}
}

package utils

import (
	"context"
	"errors"
	"github.com/a-h/templ"
	"net/http"
	"strconv"
)

var EmptyMap = map[string]string{}

func ExtractValueFromContext[T any](ctx context.Context, key string) (T, error) {
	var value T
	extractedVal := ctx.Value(key)

	if val, ok := extractedVal.(T); ok {
		value = val
		return value, nil
	}

	return value, errors.New("couldn't extract the value from the context")
}

func ConvertToTemplURL(parts ...any) templ.SafeURL {
	urlString := ""

	for _, part := range parts {
		urlString += "/"
		switch v := part.(type) {
		case int:
			urlString += strconv.Itoa(v)
		case int64:
			urlString += strconv.FormatInt(v, 10)
		case string:
			urlString += v
		default:
			// TODO somehow pass the default logger here and log the missing type check
		}
	}

	return templ.URL(urlString)
}

func ConvertToTemplStringURL(parts ...any) string {
	return string(ConvertToTemplURL(parts...))
}

type TemplateHandler struct {
	Template templ.Component
}

func (h *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := templ.Handler(h.Template)
	handler.ServeHTTP(w, r)
}

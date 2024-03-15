package utils

import (
	"context"
	"errors"
)

func ExtractValueFromContext[T any](ctx context.Context, key string) (T, error) {
	var value T
	extractedVal := ctx.Value(key)

	if val, ok := extractedVal.(T); ok {
		value = val
		return value, nil
	}

	return value, errors.New("couldn't extract the value from the context")
}

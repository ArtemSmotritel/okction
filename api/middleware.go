package api

import (
	"context"
	"fmt"
	"github.com/artemsmotritel/oktion/utils"
	"log"
	"net/http"
)

func setUserInfoToContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := extractUserIDFromCookie(r)

		if err != nil {
			r = r.WithContext(context.WithValue(r.Context(), "isAuthorized", false))
		} else {
			r = r.WithContext(context.WithValue(r.Context(), "userId", id))
			r = r.WithContext(context.WithValue(r.Context(), "isAuthorized", true))
		}

		r = r.WithContext(context.WithValue(r.Context(), "hxBoosted", r.Header.Get("HX-Boosted") == "true"))

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hxBoosted, _ := utils.ExtractValueFromContext[bool](r.Context(), "hxBoosted")

		l := fmt.Sprintf("New Request: method - %s, url - %s, hx-boosted - %t", r.Method, r.URL.Path, hxBoosted)
		logger.Println(l)

		isAuth := r.Context().Value("isAuthorized").(bool)
		l = fmt.Sprintf("User: isAuthorized - %t", isAuth)
		if isAuth {
			id := r.Context().Value("userId").(int64)
			l += fmt.Sprintf(", id - %d", id)
		}
		logger.Println(l)

		next.ServeHTTP(w, r)
	})
}

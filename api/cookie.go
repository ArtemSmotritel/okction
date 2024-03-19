package api

import (
	"net/http"
)

func getCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("userId")
	return cookie, err
}

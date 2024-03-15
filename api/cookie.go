package api

import "net/http"

func getCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("userId")
	return cookie, err
}

func setCookie(w http.ResponseWriter, _ *http.Request) {
	cookie := http.Cookie{
		Name:     "userId",
		Value:    "1",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/")
}

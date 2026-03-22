package auth

import (
	"net/http"
)

// CookieMaxAge is the lifetime of the session cookie in seconds (7 days).
const CookieMaxAge = 7 * 24 * 60 * 60

const CookieName = "gotify-client-token"

func SetCookie(w http.ResponseWriter, token string, maxAge int, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

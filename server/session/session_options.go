package session

import (
	"net/http"

	gorillasession "github.com/gorilla/sessions"
)

func NewSessionOptions() *gorillasession.Options {
	return &gorillasession.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		Secure:   false, //true
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

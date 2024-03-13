package auth

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

const (
	maxAge      = 86400
	isProd      = false
	callbackURL = "http://localhost:6060/auth/google/callback"
)

func NewAuth() {
	store := sessions.NewCookieStore([]byte("secret"))
	store.Options = &sessions.Options{
		MaxAge:   maxAge,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProd,
	}

	gothic.Store = store
}

package auth

import (
	"github.com/echo-webkom/goat/internal/db"
	"golang.org/x/oauth2"
)

type Provider interface {
	Name() string
	SetName(name string)
	BeginAuth(state string) (Session, error)
	UnmarshalSession(string) (Session, error)
	FetchUser(Session) (db.User, error)
	Debug(bool)
	RefreshToken(refreshToken string) (*oauth2.Token, error) // Get new access token based on the refresh token
	RefreshTokenAvailable() bool                             // Refresh token is provided by auth provider or not
}

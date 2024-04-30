package providers

import (
	"os"

	"github.com/echo-webkom/goat/internal/auth"
	"github.com/echo-webkom/goat/internal/domain"
	"golang.org/x/oauth2"
)

func Feide() auth.Provider {
	const (
		authUrl  = ""
		tokenUrl = ""
		userUrl  = ""

		scopeEmail   = "email"
		scopeOpenID  = "openid"
		scopeProfile = "profile"
		scopeGroups  = "groups"
	)

	getUser := func(token *oauth2.Token) (user domain.User, err error) {

		return user, err
	}

	return auth.New(
		"feide",
		os.Getenv("FEIDE_CLIENT_ID"),
		os.Getenv("FEIDE_CLIENT_SECRET"),
		authUrl,
		tokenUrl,
		[]string{
			scopeProfile,
			scopeOpenID,
			scopeEmail,
			scopeGroups,
		},
		getUser,
	)
}

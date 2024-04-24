package providers

import (
	"github.com/echo-webkom/goat/internal/auth"
	"github.com/echo-webkom/goat/internal/domain"
	"golang.org/x/oauth2"
)

func Github() auth.Provider {
	getUser := func(token *oauth2.Token) (user domain.User, err error) {

		// See:
		// https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps

		// GET https://api.github.com/user
		// Header:
		// 		"Authorization: Bearer {TOKEN_HERE}"

		return user, err
	}

	return auth.New(
		"github",
		"", // client id
		"", // client secret
		"https://github.com/login/oauth/authorize",
		"https://github.com/login/oauth/access_token",
		[]string{},
		getUser,
	)
}

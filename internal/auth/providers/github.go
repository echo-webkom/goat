package providers

import (
	"fmt"
	"os"

	"github.com/echo-webkom/goat/internal/auth"
	"github.com/echo-webkom/goat/internal/domain"
	"golang.org/x/oauth2"
)

func Github() auth.Provider {
	const (
		authUrl  = "https://github.com/login/oauth/authorize"
		tokenUrl = "https://github.com/login/oauth/authorize"
		userUrl  = "https://api.github.com/user"

		scopeUser = "user"
	)

	getUser := func(token *oauth2.Token) (user domain.User, err error) {

		// https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps

		req, err := auth.NewRequest("GET", userUrl)
		if err != nil {
			return user, err
		}

		req.AddHeader("Authorization", "Bearer "+token.AccessToken)

		body, err := req.Send()
		if err != nil {
			return user, err
		}

		fmt.Println(string(body))

		return user, err
	}

	return auth.New(
		"github",
		os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_CLIENT_SECRET"),
		authUrl,
		tokenUrl,
		[]string{
			scopeUser,
		},
		getUser,
	)
}

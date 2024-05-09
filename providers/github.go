package providers

import (
	"os"

	"github.com/echo-webkom/goat/auth"
)

func Github() auth.Provider {
	const (
		authUrl  = "https://github.com/login/oauth/authorize"
		tokenUrl = "https://github.com/login/oauth/access_token"

		scopeUser = "user"
	)

	return auth.New(
		"github",
		os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_CLIENT_SECRET"),
		authUrl,
		tokenUrl,
		[]string{
			scopeUser,
		},
	)
}

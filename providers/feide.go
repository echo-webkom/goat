package providers

import (
	"os"

	"github.com/echo-webkom/goat/auth"
)

func Feide() auth.Provider {
	const (
		authUrl  = "https://auth.dataporten.no/oauth/authorization"
		tokenUrl = "https://auth.dataporten.no/oauth/token"
		userUrl  = "https://auth.dataporten.no/openid/userinfo"

		scopeEmail   = "email"
		scopeOpenID  = "openid"
		scopeProfile = "profile"
		scopeGroups  = "groups"
	)

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
	)
}

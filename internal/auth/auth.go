package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/echo-webkom/goat/internal/domain"
	"golang.org/x/oauth2"
)

type UserFetcher func(*oauth2.Token) (domain.User, error)

type Provider struct {
	name   string
	config oauth2.Config

	// Given as an argument when calling auth.New(). getuser should
	// fetch user data from the provider using the given token and
	// return a User struct.
	getUser UserFetcher
}

// New creates a new provider with oauth config
func New(
	providerName,
	clientId,
	clientSecret,
	authUrl,
	tokenUrl string,
	scopes []string,
	getUser UserFetcher,
) Provider {

	callbackUrl := "http://localhost:8080/auth/" + providerName + "/callback"
	return Provider{
		config: oauth2.Config{
			RedirectURL:  callbackUrl,
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Scopes:       []string{},
			Endpoint: oauth2.Endpoint{
				AuthURL:  authUrl,
				TokenURL: tokenUrl,
			},
		},
		name:    providerName,
		getUser: getUser,
	}
}

// BeginAuthHandler handles the endpoint /auth/{provider}
// Redirects the user to the providers auth page.
func BeginAuthHandler(providers map[string]Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerName := r.PathValue("provider")

		if p, ok := providers[providerName]; ok {
			// Todo: set state cookie
			// Todo: add random state
			url := p.config.AuthCodeURL("abcdefppabcdefppabcdefppabcdefpp")
			http.Redirect(w, r, url, http.StatusSeeOther)

		} else {
			log.Println("login: provider not found: ", providerName)
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

// CallbackHandler handles the endpoint /auth/{provider}/callback.
func CallbackHandler(providers map[string]Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := providers[r.PathValue("provider")]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Todo: compare state

		code := r.FormValue("code")
		token, err := p.config.Exchange(context.Background(), code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := p.getUser(token)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		AsAuthenticatedUser(user)
	}
}

func AsAuthenticatedUser(user domain.User) {
	// ... do stuff with user
}

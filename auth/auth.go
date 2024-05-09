package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type UserFetcher func(*oauth2.Token) error

type Provider struct {
	name   string
	config oauth2.Config
}

// New creates a new provider with oauth config
func New(
	providerName,
	clientId,
	clientSecret,
	authUrl,
	tokenUrl string,
	scopes []string,
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
		name: providerName,
	}
}

// BeginAuthHandler handles the endpoint /auth/{provider}
// Redirects the user to the providers auth page.
func BeginAuthHandler(providers map[string]Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerName := r.PathValue("provider")

		if p, ok := providers[providerName]; ok {
			b := make([]byte, 32+2)
			rand.Read(b)
			state := fmt.Sprintf("%x", b)[2 : 32+2]

			http.SetCookie(w, &http.Cookie{
				Name:     "state",
				Value:    state,
				Secure:   true,
				HttpOnly: true,
				Path:     "/",
				SameSite: http.SameSiteLaxMode,
			})

			url := p.config.AuthCodeURL(state)
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
			log.Println("callback: provider not found")
			return
		}

		cookie, err := r.Cookie("state")
		if err != nil {
			log.Println("callback: state cookie not found")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.URL.Query().Get("state") != cookie.Value {
			log.Println("callback: state cookie mismatch")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Welcome to "guess how the provider gives you the code"
		guesses := []string{
			r.FormValue("code"),
			r.URL.Query().Get("code"),
		}

		code := ""
		for _, c := range guesses {
			if c != "" {
				code = c
				break
			}
		}

		if code == "" {
			log.Println("callback: couldnt get code")
		}

		token, err := p.config.Exchange(context.Background(), code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			log.Println("callback: token exchange failed")
			return
		}

		AsAuthenticatedUser(token)
	}
}

func AsAuthenticatedUser(token *oauth2.Token) {
	log.Println("authenticated as user")
	// ... do stuff with user
}

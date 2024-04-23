package auth

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type Provider struct {
	config oauth2.Config
}

func New(clientId, clientSecret, authUrl, tokenUrl string) Provider {
	config := oauth2.Config{
		// Default auth callback for testing. Remove
		RedirectURL:  "http://localhost:8080/auth/{provider}/callback",
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authUrl,
			TokenURL: tokenUrl,
		},
	}

	return Provider{config: config}
}

type User struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	Error        string `json:"error"`
	ErrorDesc    string `json:"error_description"`
	ErrorUri     string `json:"error_uri"`
}

type Session struct {
	User    User
	Writer  http.ResponseWriter
	Request *http.Request
}

// LoginHandler return a http.HandlerFunc used to handle the endpoint /auth/{provider}
func LoginHandler(providers map[string]Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerName := r.PathValue("provider")
		if p, ok := providers[providerName]; ok {
			url := p.config.AuthCodeURL("randomstate")
			http.Redirect(w, r, url, http.StatusSeeOther)
		} else {
			log.Println("login: provider not found: ", providerName)
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

// CallbackHandler returns a http.HandlerFunc used to handle the endpoint /auth/{provider}/callback
func CallbackHandler(providers map[string]Provider, userCallback func(Session)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerName := r.PathValue("provider")

		if p, ok := providers[providerName]; ok {
			userData, err := p.fetchUserData(r)
			if err != nil {
				log.Println(err)
				return
			}

			var user User
			if err = json.Unmarshal(userData, &user); err != nil {
				log.Println(err)
				return
			}

			userCallback(Session{
				User:    user,
				Request: r,
				Writer:  w,
			})
		} else {
			log.Println("callback: provider not found: ", providerName)
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func (p *Provider) fetchUserData(r *http.Request) ([]byte, error) {
	code := r.FormValue("code")

	token, err := p.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	token_url := p.config.Endpoint.TokenURL + "?access_token=" + token.AccessToken
	resp, err := http.Get(token_url)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

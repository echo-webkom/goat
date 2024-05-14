package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

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
			// Make random state string
			b := make([]byte, 32+2)
			rand.Read(b)
			state := fmt.Sprintf("%x", b)[2 : 32+2]

			setSecureCookie(w, "state", state)
			// Todo: use other method maybe?
			setSecureCookie(w, "origin", r.URL.Query().Get("origin"))

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

		if err := compareState(r); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		code, err := getCode(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		token, err := p.config.Exchange(context.Background(), code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("callback: token exchange failed: %s\n", err.Error())
			return
		}

		AsAuthenticatedUser(w, r, p, token)
	}
}

func AsAuthenticatedUser(w http.ResponseWriter, r *http.Request, provider Provider, token *oauth2.Token) {
	log.Println("authenticated as user")

	cookie, err := r.Cookie("origin")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	signingKey := []byte(os.Getenv("JWT_KEY"))

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token":    token.AccessToken,
		"provider": provider.name,
	})

	str, err := tok.SignedString(signingKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	url := cookie.Value + "?jwt=" + str
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func setSecureCookie(w http.ResponseWriter, key, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    value,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
}

func compareState(r *http.Request) error {
	cookie, err := r.Cookie("state")
	if err != nil {
		return errors.New("state cookie not found")
	}

	if r.URL.Query().Get("state") != cookie.Value {
		return errors.New("state mismatch")
	}

	return nil
}

func getCode(r *http.Request) (string, error) {
	guesses := []string{
		r.FormValue("code"),
		r.URL.Query().Get("code"),
	}

	for _, c := range guesses {
		if c != "" {
			return c, nil
		}
	}

	return "", errors.New("could not extract code from request")
}

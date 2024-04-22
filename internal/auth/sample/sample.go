package sample

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

func resJson(w http.ResponseWriter, j any) {
	b, err := json.Marshal(j)
	if err != nil {
		log.Println("JSON marshal error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

// Mount endpoints handled by provider for testing
func mountExampleHandlers(s *http.ServeMux) {
	// Example login page, will be replaced with provider URL
	s.HandleFunc("GET /sample/auth", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/auth/sample/sample_auth.html")
	})

	// Used for token exchange
	s.HandleFunc("POST /sample/tokenUrl", func(w http.ResponseWriter, r *http.Request) {
		resJson(w, map[string]any{
			"access_token":  "abcdef",
			"token_type":    "bearer",
			"expires_in":    3600,
			"refresh_token": "ghijklmno",
			"scope":         "",
		})
	})

	// Used to fetch user data with generated token
	s.HandleFunc("GET /sample/tokenUrl", func(w http.ResponseWriter, r *http.Request) {
		resJson(w, map[string]any{
			"username":     "bob",
			"access_token": r.URL.Query().Get("access_token"),
		})
	})
}

// Todo: create generic newProvider function

func New(s *http.ServeMux) {
	mountExampleHandlers(s)

	const (
		// Load from .env
		CLIENT_ID     = "john"
		CLIENT_SECRET = "1234"

		AUTH_URL  = "http://localhost:8080/sample/auth"
		TOKEN_URL = "http://localhost:8080/sample/tokenUrl"
	)

	config := oauth2.Config{
		RedirectURL:  "http://localhost:8080/sample_callback",
		ClientID:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  AUTH_URL,
			TokenURL: TOKEN_URL,
		},
	}

	s.Handle("GET /sample_login", login(config))
	s.Handle("POST /sample_callback", callback(config))
}

// Creates new login handler. Should redirect to providers auth URL with
// generated state. URL is given client id/secret, redirect uri, callback
// uri and state.
func login(config oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := config.AuthCodeURL("randomstate")
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

// Creates a new callback handler for the auth provider. Verifies state
// and creates access token from code given by provider. This handler simply
// responds with the user data json.
func callback(config oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")
		if state != "randomstate" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		code := r.FormValue("code")

		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		token_url := config.Endpoint.TokenURL + "?access_token=" + token.AccessToken
		resp, err := http.Get(token_url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userData, err := io.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(userData)
	}
}

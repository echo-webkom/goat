package sample

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/echo-webkom/goat/internal/auth"
	"github.com/echo-webkom/goat/internal/domain"
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
func MountExampleHandlers(s *http.ServeMux) {
	// Example login page, will be replaced with provider URL
	s.HandleFunc("GET /sample/auth", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/auth/sample/sample_auth.html")
	})

	// Used for token exchange
	s.HandleFunc("POST /sample/tokenUrl", func(w http.ResponseWriter, r *http.Request) {
		resJson(w, map[string]any{
			"access_token":  "VeryCoolAccessToken",
			"token_type":    "bearer",
			"expires_in":    3600,
			"refresh_token": "CoolerRefreshToken",
			"scope":         "CoolSCope",
		})
	})

	// Used to fetch user data with generated token
	s.HandleFunc("GET /sample/tokenUrl", func(w http.ResponseWriter, r *http.Request) {
		resJson(w, map[string]any{
			"access_token": r.URL.Query().Get("access_token"),
		})
	})
}

func New() auth.Provider {
	return auth.New(
		"sample",
		"client_id",
		"client_secret",
		"http://localhost:8080/sample/auth",
		"http://localhost:8080/sample/tokenUrl",
		[]string{},
		func(t *oauth2.Token) (u domain.User, e error) {
			return u, e
		},
	)
}

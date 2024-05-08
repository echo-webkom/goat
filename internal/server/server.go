package server

import (
	"net/http"

	"github.com/echo-webkom/goat/internal/auth"
	"github.com/echo-webkom/goat/internal/auth/providers"
)

type Server struct {
	Router *http.ServeMux
	Config Config
}

type Config struct {
	Addr string
}

func New() *Server {
	return &Server{
		Router: http.NewServeMux(),
		Config: Config{
			Addr: ":8080",
		},
	}
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.Router)
}

func (s *Server) MountHandlers() {
	ps := map[string]auth.Provider{
		"github": providers.Github(),
		"feide":  providers.Feide(),
	}

	s.Router.HandleFunc("/auth/{provider}", auth.BeginAuthHandler(ps))
	s.Router.HandleFunc("/auth/{provider}/callback", auth.CallbackHandler(ps))
}

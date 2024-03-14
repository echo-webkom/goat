package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Config struct {
	Addr string
}

type Server struct {
	Router *chi.Mux
	Config *Config
}

func New() *Server {
	r := chi.NewRouter()

	cfg := &Config{
		Addr: ":8080",
	}

	return &Server{
		Router: r,
		Config: cfg,
	}
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.Router)
}

func (s *Server) MountHandlers() {

	handler := NewHandler(func(hwc HandlerWithCtx) {
		hwc.w.Write([]byte("Hello " + hwc.name))
	})

	middleware := NewMiddleware(func(hf HandlerFunc) HandlerFunc {
		return func(hwc HandlerWithCtx) {
			hwc.name = "John"
			hf(hwc)
		}
	})

	s.Router.Get("/", ToHttpHandler(middleware(handler)))
}

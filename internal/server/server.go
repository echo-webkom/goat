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

	// Create simple base handler using a context
	handler := NewHandler(func(hwc HandlerWithCtx) {
		hwc.res.Write([]byte("Hello " + hwc.name))
	})

	// Create middleware that writes to the context before calling the handler
	middleware := NewMiddleware(func(hf HandlerFunc) HandlerFunc {
		return func(hwc HandlerWithCtx) {
			hwc.name = "John"
			hf(hwc)
		}
	})

	s.Router.Get("/", ToHttpHandler(middleware(handler)))
}

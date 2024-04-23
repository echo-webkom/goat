package server

import (
	"encoding/json"
	"net/http"

	"github.com/echo-webkom/goat/internal/auth"
	"github.com/echo-webkom/goat/internal/auth/sample"
)

type Server struct {
	Router *http.ServeMux
	Config Config
}

type Config struct {
	Addr string
}

func New() *Server {
	r := http.NewServeMux()

	cfg := Config{
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
	handler := NewHandler(func(ctx Context) {
		ctx.res.Write([]byte("Hello " + ctx.name))
	})

	// Create myMiddleware that writes to the context before calling the handler
	middleware := NewMiddleware(func(hf HandlerFunc) HandlerFunc {
		return func(ctx Context) {
			ctx.name = "John"
			hf(ctx)
		}
	})

	s.Router.Handle("GET /test", ToHttpHandlerFunc(middleware(handler)))

	// Sample oauth2 flow, go to /auth/sample
	ps := map[string]auth.Provider{
		"sample": sample.New(),
	}

	s.Router.HandleFunc("/auth/{provider}", auth.LoginHandler(ps))
	s.Router.HandleFunc("/auth/{provider}/callback", auth.CallbackHandler(ps, func(s auth.Session) {
		d, _ := json.Marshal(s.User)
		s.Writer.Write(d)
	}))

	sample.MountExampleHandlers(s.Router)
}

package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/echo-webkom/goat/auth"
	"github.com/echo-webkom/goat/providers"
)

func New() *http.Server {
	mux := http.NewServeMux()

	ps := map[string]auth.Provider{
		"github": providers.Github(),
		"feide":  providers.Feide(),
	}

	mux.HandleFunc("/auth/{provider}", auth.BeginAuthHandler(ps))
	mux.HandleFunc("/auth/{provider}/callback", auth.CallbackHandler(ps))

	// Test
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "temp/auth.html")
	})
	mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "temp/home.html")
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func ServeWithShutdown(s *http.Server) {
	log.Println("Listening at", s.Addr)

	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections. Shutting down...")
	}()

	// Handle interupt (ctrl-c etc) and gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	log.Println("Graceful shutdown complete.")
}

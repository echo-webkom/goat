package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type RequestWrapper func(http.Handler) http.Handler

type goatHandler struct {
	h  http.HandlerFunc
	id string
}

func (h goatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.h.ServeHTTP(w, r)
}

func newHandler(h http.HandlerFunc) goatHandler {
	return goatHandler{
		func(w http.ResponseWriter, r *http.Request) {

		},
		"1234",
	}
}

func handleHellWorld(name string) goatHandler {
	msg := "Hello " + name

	helloWorld := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(msg))

	})

	checkAdmin := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// if !isAdmin(r) {
			// 	http.NotFound(w, r)
			// }

			h.ServeHTTP(w, r)
		})
	}

	return checkAdmin(helloWorld)
}

func someMiddelWare(handler goatHandler) goatHandler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("new request!")
		handler.ServeHTTP(w, r)
	})
}

func NewServer() http.Handler {
	mux := chi.NewMux()

	mux.Handle("/", handleHellWorld("john"))

	var handler http.Handler = mux
	handler = someMiddelWare(handler)

	return handler
}

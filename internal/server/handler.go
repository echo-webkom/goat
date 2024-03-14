package server

import "net/http"

type HandlerWithCtx struct {
	res  http.ResponseWriter
	req  *http.Request
	name string
}

type HandlerFunc func(HandlerWithCtx)

func NewHandler(f HandlerFunc) HandlerFunc {
	return f
}

func ToHttpHandler(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(HandlerWithCtx{w, r, ""}) // Create empty base context
	}
}

type Middleware func(HandlerFunc) HandlerFunc

func NewMiddleware(m Middleware) Middleware {
	return m
}

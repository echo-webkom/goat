package server

import "net/http"

type Context struct {
	res http.ResponseWriter
	req *http.Request

	name string
}

type HandlerFunc func(Context)

func NewHandler(f HandlerFunc) HandlerFunc {
	return f
}

func ToHttpHandlerFunc(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(Context{w, r, ""}) // Create empty base context
	}
}

func ToHttpHandler(f HandlerFunc) http.Handler {
	return ToHttpHandlerFunc(f)
}

type Middleware func(HandlerFunc) HandlerFunc

func NewMiddleware(m Middleware) Middleware {
	return m
}

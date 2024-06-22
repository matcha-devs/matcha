package main

import (
	"net/http"
	"time"
)

const maxHandleTime = time.Second

func withClientTimeout(handlerFunc http.HandlerFunc) http.Handler {
	return http.TimeoutHandler(handlerFunc, maxHandleTime, "")
}

func router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /public/", withClientTimeout(servePublicFile))
	mux.Handle("GET /{$}", withClientTimeout(loadIndex))
	mux.Handle("POST /signup-submit", withClientTimeout(signupSubmit))
	mux.Handle("POST /login-submit", withClientTimeout(loginSubmit))
	mux.Handle("POST /delete-user", withClientTimeout(deleteUser))
	mux.Handle("GET /", withClientTimeout(loadPage))
	return mux
}

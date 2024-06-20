package main

import (
	"net/http"
	"time"
)

var (
	maxHandleTime    = time.Second
	validEntryPoints = map[string]struct{}{
		"signup": {}, "signup-submit": {}, "signup-fail": {},
		"login": {}, "login-submit": {}, "login-fail": {},
		"dashboard": {}, "settings": {}, "delete-user": {},
	}
)

func withClientTimeout(handlerFunc http.HandlerFunc) http.Handler {
	return http.TimeoutHandler(handlerFunc, maxHandleTime, "")
}

func router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", withClientTimeout(loadEntryPoint))
	mux.Handle("GET /{$}", withClientTimeout(loadIndex))
	mux.Handle("GET /public/", withClientTimeout(servePublicFile))
	mux.Handle("POST /signup-submit", withClientTimeout(signupSubmit))
	mux.Handle("POST /login-submit", withClientTimeout(loginSubmit))
	mux.Handle("POST /delete-user", withClientTimeout(deleteUser))
	return mux
}

package main

import (
	"log"
	"net/http"
	"time"
)

const maxHandleTime = time.Second

func withRequestLogs(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("Request:", r.Method, r.URL)
			handler.ServeHTTP(w, r)
		},
	)
}

func withClientTimeout(handlerFunc http.HandlerFunc) http.Handler {
	return http.TimeoutHandler(handlerFunc, maxHandleTime, "")
}

func loggedRouter() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /public/", withClientTimeout(servePublicFile))
	mux.Handle("GET /{$}", withClientTimeout(loadIndex))
	mux.Handle("POST /signup", withClientTimeout(signupSubmit))
	mux.Handle("POST /login", withClientTimeout(loginSubmit))
	mux.Handle("POST /delete-user", withClientTimeout(deleteUser))
	mux.Handle("GET /", withClientTimeout(loadPage))
	return withRequestLogs(mux)
}

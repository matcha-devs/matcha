package main

import (
	"log"
	"net/http"
	"time"
)

const maxHandleTime = 5 * time.Second

func withRequestLogs(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request:", r.Method, r.URL)
		handler.ServeHTTP(w, r)
	}
}

func withClientTimeout(handlerFunc http.HandlerFunc) http.Handler {
	return http.TimeoutHandler(handlerFunc, maxHandleTime, "")
}

func loggedRouter() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /public/", withClientTimeout(getPublic))
	mux.Handle("GET /{$}", withClientTimeout(getIndex))
	mux.Handle("POST /signup", withClientTimeout(postSignup))
	mux.Handle("POST /login", withClientTimeout(postLogin))
	mux.Handle("POST /logout", withClientTimeout(postLogout))
	mux.Handle("POST /delete-user", withClientTimeout(postDeleteUser))
	mux.Handle("GET /", withClientTimeout(getPage))
	return withRequestLogs(mux)
}

package main

import (
	"log"
	"net/http"
	"strings"
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

func loadIndex(w http.ResponseWriter, r *http.Request) {
	loadPage(w, r, "index")
}

func entryPoint(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	log.Println("Routing {" + path + "}")
	if _, exists := validEntryPoints[path]; !exists {
		log.Println("Not a valid entry point:", path)
		http.NotFound(w, r)
	}
	loadPage(w, r, path)
}

func handlerWithTimeout(handlerFunc http.HandlerFunc) http.Handler {
	return http.TimeoutHandler(handlerFunc, maxHandleTime, "")
}

func newMatchaRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /{$}", handlerWithTimeout(loadIndex))
	mux.Handle("POST /signup-submit", handlerWithTimeout(signupSubmit))
	mux.Handle("POST /login-submit", handlerWithTimeout(loginSubmit))
	mux.Handle("POST /delete-user", handlerWithTimeout(deleteUser))
	mux.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	mux.Handle("/", handlerWithTimeout(entryPoint)) // TODO(@CarlosACJ55): Improve handling of entry points.
	return mux
}

// TODO(@FaaizMemonPurdue): This is an example of how go routines should be used, but we still need API call timeouts
// func routeWithTimeout(w http.ResponseWriter, r *http.Request) {
// 	ctx, cancel := context.WithTimeout(r.Context(), maxHandleTime)
// 	defer cancel()
// 	select {
// 	case <-ctx.Done():
// 		log.Println("Routing took longer than", maxHandleTime)
// 	default:
// 		start := time.Now()
// 		entryPoint(w, r)
// 		log.Println("Routing done after", time.Since(start))
// 	}
// }

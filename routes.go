package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

func route(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	log.Println("Routing {" + path + "}")
	switch path {
	case "":
		loadPage(w, r, "index")
	case "signup-submit":
		signupSubmit(w, r)
	case "login-submit":
		loginSubmit(w, r)
	case "delete-user":
		deleteUser(w, r)
	default:
		if _, exists := validEntryPoints[path]; exists {
			loadPage(w, r, path)
		} else {
			http.NotFound(w, r)
		}
	}
}

// TODO(@FaaizMemonPurdue): This is an example of how go routines should be used, but we still need server timeouts
func routeWithTimeout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), maxRouteTime)
	defer cancel()
	select {
	case <-ctx.Done():
		log.Println("Routing took longer than", maxRouteTime)
	//case <-time.After(maxRouteTime):
	default:
		start := time.Now()
		route(w, r)
		log.Println("Routing done after", time.Since(start))
	}
}

func matchaMux() *http.ServeMux {
	mux := http.NewServeMux()
	// TODO(@CarlosACJ55): Make a clean transition from the switch case to ServeMux
	//mux.Handle("/{$}", http.TimeoutHandler(http.HandlerFunc(loadPage), maxRouteTime, ""))
	mux.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	mux.Handle("/", http.TimeoutHandler(http.HandlerFunc(routeWithTimeout), maxRouteTime, ""))
	return mux
}

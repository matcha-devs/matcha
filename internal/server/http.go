package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type HTTPServer struct {
	underlyingServer http.Server
}

func New(handler http.Handler) *HTTPServer {
	return &HTTPServer{
		http.Server{
			Addr:                         ":8080",
			Handler:                      handler,
			DisableGeneralOptionsHandler: false,
			TLSConfig:                    nil,
			ReadTimeout:                  time.Second,
			ReadHeaderTimeout:            2 * time.Second,
			WriteTimeout:                 time.Second,
			IdleTimeout:                  30 * time.Second,
			MaxHeaderBytes:               0,
			TLSNextProto:                 nil,
			ConnState:                    nil,
			ErrorLog:                     nil,
			BaseContext:                  nil,
			ConnContext:                  nil,
		},
	}
}

func (s *HTTPServer) Run() (err error) {
	log.Println("HTTP server starting on", s.underlyingServer.Addr, "🫡")
	if err = s.underlyingServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln("HTTP server run error -", err)
	}
	return
}

func (s *HTTPServer) Shutdown(maxClientDisconnectTime time.Duration) (err error) {
	ctx, release := context.WithTimeout(context.Background(), maxClientDisconnectTime)
	defer release()
	if err = s.underlyingServer.Shutdown(ctx); err != nil {
		log.Println("HTTP server close error -", err)
	} else {
		log.Println("HTTP server has shutdown 👋🏽")
	}
	return
}

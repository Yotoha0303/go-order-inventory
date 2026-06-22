package app

import (
	"fmt"
	"net/http"
)

func NewHTTPServer(deps *Deps) *http.Server {
	cfg := deps.Config.HttpServer.Server

	return &http.Server{
		Addr:              fmt.Sprintf(":%d", deps.Config.Server.Port),
		Handler:           deps.Router,
		ReadTimeout:       cfg.ReadTimeOut,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytesKib << 10,
	}
}

package rest

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	readTimeout    = 10 * time.Second
	writeTimeout   = 30 * time.Second
	idleTimeout    = 0 * time.Second
	maxHeaderBytes = 1 << 20
)

func (api *API) ListenAndServe(servicePort string) {
	server := http.Server{
		Addr:           fmt.Sprintf(":%s", servicePort),
		Handler:        api.router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		IdleTimeout:    idleTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	api.server = &server

	go func() {
		ln, err := net.Listen("tcp", server.Addr)
		if err != nil {
			log.Fatalf("Unable to listen on %s: %s", server.Addr, err.Error())
		}

		if err := server.Serve(ln); err != http.ErrServerClosed {
			log.Fatalf("Failed serve connections: %s", err.Error())
		}
	}()

	log.Printf("Listening on %s.", server.Addr)
}

func (api *API) Shutdown(ctx context.Context) {
	if err := api.server.Shutdown(ctx); err != nil {
		log.Printf("Unable to shutdown server: %s", err.Error())
	}

	log.Printf("Server gracefully stopped listening on %s.", api.server.Addr)
}

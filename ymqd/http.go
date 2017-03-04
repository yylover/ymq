package ymqd

import (
	"fmt"
	"time"
)

type HttpServer struct {
	server *net.Server
}

func NewHTTPServer() (*HttpServer, error) {
	server := new(HttpServer)

	server.server = &Http.Server{
		Addr:           ":8082",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return server, nil
}

func (h *HttpServer) init() error {
	h.server.ListenAndServe()
	fmt.Print("a")
}
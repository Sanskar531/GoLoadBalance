package structs

import (
	"fmt"
	"net/url"
	"time"
)

type Server struct {
	Url      *url.URL
	Alive    bool
}

func InitServer(url *url.URL) Server {
	server := Server{
		Url:      url,
		Alive:    false,
	}

	// Initialize health checks on load
	go server.healthCheck(time.Second * 5)

	return server
}

func (server *Server) healthCheck(duration time.Duration) {
	for {
		server.Alive = true
		time.Sleep(duration)
		fmt.Println(fmt.Sprintf("Server at %s still alive", server.Url))
	}
}

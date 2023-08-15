package structs

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Server struct {
	Url   *url.URL
	Alive bool
}

func InitServer(url *url.URL) *Server {
	server := Server{
		Url:   url,
		Alive: true,
	}

	// Initialize health checks on load
	go server.healthCheck(time.Second * 5)

	return &server
}

func (server *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Forwarding Request to path %s", r.URL.Path);
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write([]byte("<h1>Hello!</h1>"))
	if err != nil {
		log.Fatal(err)
	}
}

func pingServer(url string) bool {
	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{
		Transport: transportConfig,
	}

	_, err := client.Get(url)

	if err != nil {
		// Log directly prints to std err
		log.Printf("Server is not alive at %s", url)
		return false
	} else {
		return true
	}
}

// Ping the host at the /ping route and check if the server
// is still responding
// TODO: get a function as an argument to handle the case where
// the server is not responding
func (server *Server) healthCheck(duration time.Duration) {
	for {
		pingUrl, err := url.JoinPath(server.Url.String(), "ping")

		if err != nil {
			log.Fatal(err)
		}

		server.Alive = pingServer(pingUrl)
		time.Sleep(duration)
	}
}

package lib

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

type Server struct {
	Url               *url.URL
	Alive             bool
	client            *http.Client
	retryCount        int
	isDeadChannel     *(chan bool)
	ActiveConnections int32
}

func InitServer(url *url.URL, healthCheckFrequencyInSeconds int, maxRetryCount int) *Server {
	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{
		Transport: transportConfig,
	}

	isDeadChannel := make(chan bool)
	server := Server{
		Url: url,
		// We initialize as false as in the start we don't know if the host is alive
		Alive:         false,
		client:        &client,
		isDeadChannel: &isDeadChannel,
	}

	// Initialize health checks on load
	go server.healthCheck(time.Second*time.Duration(healthCheckFrequencyInSeconds), maxRetryCount)
	return &server
}

func (server *Server) HandleRequest(responseWriter http.ResponseWriter, request *http.Request) (*http.Response, *string, error) {
	log.Printf("Forwarding Request to path %s using host %s", request.URL.Path, server.Url.Host)
	// We don't need to keep the old host intact as that host is us
	// We remove the old host and reuse the same request by changing the host to the server.
	request.URL.Host = server.Url.Host
	request.URL.Scheme = server.Url.Scheme
	// RequestURI needs to be empty for this to be a client request.
	request.RequestURI = ""

	// Multiple threads will access this
	atomic.AddInt32(&server.ActiveConnections, 1)

	res, err := server.client.Do(request)

	// Multiple threads will access this
	atomic.AddInt32(&server.ActiveConnections, -1)

	if err != nil {
		log.Println("Error: ", err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return nil, nil, err
	}

	// Reset the Headers properly before relaying the response back
	responseWriter.Header().Set("Content-Length", res.Header.Get("Content-Length"))
	responseWriter.Header().Set("Content-Type", res.Header.Get("Content-Type"))
	responseWriter.Header().Set("Status", res.Status)
	// Body is a stream we read from the stream into two and write it to the reader we get back
	// from tee and then write it to the responseWriter as well.
	teeReader := io.TeeReader(res.Body, responseWriter)
	body, err := io.ReadAll(teeReader)

	if err != nil {
		log.Println("Error while trying to copy body for caching")
	}

	res.Body.Close()
	stringifiedBody := string(body)
	return res, &stringifiedBody, nil
}

func (server *Server) ping() bool {
	pingUrl, err := url.JoinPath(server.Url.String(), "ping")

	if err != nil {
		log.Printf("Server couldn't be ping as the /ping couldn't be joined to host's url. HOST: %s", server.Url.Host)
	}

	_, err = server.client.Get(pingUrl)

	if err != nil {
		// Log directly prints to std err
		log.Printf("Server is not alive at %s", server.Url.Host)
		return false
	} else {
		return true
	}
}

// Ping the host at the /ping route and check if the server
// is still responding
// TODO: get a function as an argument to handle the case where
// the server is not responding
func (server *Server) healthCheck(duration time.Duration, maxRetryCount int) {
	for {
		if server.retryCount == maxRetryCount {
			log.Printf("Removing server at host: %s after %d retries", server.Url.Host, maxRetryCount)

			// launch webhook
			if false {
				log.Printf("Sending webhook event to %s because %s died.", "sdf", "dsf")
			}
			close(*server.isDeadChannel)

			return
		}
		isAlive := server.ping()
		// no lock required here because this is an atomic operation
		server.Alive = isAlive
		if !isAlive {
			server.retryCount += 1
		} else {
			server.retryCount = 0
		}
		time.Sleep(duration)
	}
}

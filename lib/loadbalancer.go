package lib

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sanskar531/goloadbalance/configuration"
)

type LoadBalancer struct {
	Servers               []*Server
	balancer              *Balancer
	config                *configuration.Config
	cache                 *Cache
	cacheTimeoutInSeconds int
	mutex                 *sync.RWMutex
}

type AddServerRequest struct {
	Host string `json:"host"`
}

type RemoveServerRequest struct {
	Host string `json:"host"`
}

func InitLoadBalancer(servers []*Server, balancer *Balancer, config *configuration.Config) *LoadBalancer {
	loadbalancer := &LoadBalancer{
		Servers:  servers,
		balancer: balancer,
		mutex:    &sync.RWMutex{},
		config:   config,
	}

	if config.CacheEnabled {
		loadbalancer.cache = InitCache()
		loadbalancer.cacheTimeoutInSeconds = config.CacheTimeoutInSeconds
	}

	for _, server := range servers {
		// Listen for dead servers on a different thread
		go loadbalancer.listenForDeadServer(server)
	}

	return loadbalancer
}

func (loadBalancer *LoadBalancer) getServerToHandleRequest() *Server {
	loadBalancer.mutex.RLock()
	defer loadBalancer.mutex.RUnlock()

	return (*loadBalancer.balancer).GetServer(loadBalancer.Servers)
}

// This function is called as a go routine by the http module
// when serving a request
func (loadBalancer *LoadBalancer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	// Don't process request from blocked ips
	clientIp := net.ParseIP(strings.Split(request.RemoteAddr, ":")[0])
	log.Printf(clientIp.String())
	for _, ip := range loadBalancer.config.BlacklistedIps {
		if ip.Equal(clientIp) {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			responseWriter.Write([]byte("Unauthorized request."))
			return
		}
	}

	if loadBalancer.cache != nil {
		if cachedMap := loadBalancer.cache.check(request); cachedMap != nil {
			cachedResponse := (*cachedMap)["response"].(*http.Response)
			cachedBody := (*cachedMap)["body"].(*string)

			if cachedResponse != nil {
				log.Printf("Cache Hit: Found for route %s", request.URL.Path)

				// Reset the Headers properly before relaying the response back
				responseWriter.Header().Set("Content-Length", cachedResponse.Header.Get("Content-Length"))
				responseWriter.Header().Set("Content-Type", cachedResponse.Header.Get("Content-Type"))
				responseWriter.Header().Set("Status", cachedResponse.Status)
				io.WriteString(responseWriter, *cachedBody)
				cachedResponse.Body.Close()
				return
			}
		}
	}

	server := loadBalancer.getServerToHandleRequest()
	res, body, err := server.HandleRequest(responseWriter, request)

	if loadBalancer.cache != nil && err == nil {
		// Handle caching on a different thread so that we can return the response
		go loadBalancer.cache.save(request, body, res, loadBalancer.cacheTimeoutInSeconds)
	}
}

func (loadBalancer *LoadBalancer) GetServersStatus() map[string]bool {
	serverStatuses := make(map[string]bool)
	for _, server := range loadBalancer.Servers {
		server.mutex.RLock()
		serverStatuses[server.Url.Host] = server.Alive
		server.mutex.RUnlock()
	}
	return serverStatuses
}

func (loadBalancer *LoadBalancer) listenForDeadServer(server *Server) {
	for {
		select {
		case <-*server.isDeadChannel:
			loadBalancer.gracefullyShutdownServer(server)
			// Once the server is shutdown we release the thread
			return
		}
	}
}

func (loadBalancer *LoadBalancer) FindServer(url *url.URL) (*Server, error) {
	loadBalancer.mutex.RLock()
	defer loadBalancer.mutex.RUnlock()

	for _, server := range loadBalancer.Servers {
		if server.Url.Host == url.Host {
			return server, nil
		}
	}
	return nil, &NoServerFoundErr{}
}

func (loadBalancer *LoadBalancer) AddServer(host string) error {
	parsedHostUrl, err := url.Parse(host)

	if err != nil {
		log.Printf("Can't parse host url %s", host)
	}

	_, err = loadBalancer.FindServer(parsedHostUrl)

	if err == nil {
		return &ServerExists{}
	}

	server := InitServer(parsedHostUrl, loadBalancer.config)

	loadBalancer.mutex.Lock()
	defer loadBalancer.mutex.Unlock()

	loadBalancer.Servers = append(loadBalancer.Servers, server)

	// Make sure we let our balancer know that we have a new server
	(*loadBalancer.balancer).ServerAdd()

	return nil
}

func (loadBalancer *LoadBalancer) gracefullyShutdownServer(deadServer *Server) {
	loadBalancer.mutex.Lock()
	defer loadBalancer.mutex.Unlock()

	// First, we let our balancing algorithm know that a server has died and it should
	// clean up.
	(*loadBalancer.balancer).ServerDead()

	// Second, We take it out of the server pool. For this we need to lock the resource.
	// This shoulnd't affect performance as I would assume this shouldn't happen often
	for idx, server := range loadBalancer.Servers {
		// Check if mem addresses match
		if server == deadServer {
			// If the only active server dies, we quit as well.
			if len(loadBalancer.Servers) == 1 {
				log.Fatal("Died because no other servers alive.")
			}

			// Replace the current index with the mem address for the last server
			loadBalancer.Servers[idx] = loadBalancer.Servers[len(loadBalancer.Servers)-1]
			// Remove the copied last server mem address from the server
			loadBalancer.Servers = loadBalancer.Servers[:len(loadBalancer.Servers)-1]
		}
	}

	// Once taken out of the server pool. We make sure that all it's active connections
	// finish before we clear out the resource and since this is already running on another
	// thread we can just poll and wait for the active connections to go down to 0.
	for atomic.LoadInt32(&deadServer.ActiveConnections) != 0 {
		time.Sleep(time.Second)
	}

	// Now the server should go out of scope and cleaned up by the GC.
}

func (loadBalancer *LoadBalancer) Balance() {
	// A catch all which is needed for the load balancer to redirect
	http.Handle(
		"/",
		loadBalancer,
	)

	// Utility Handler to check which hosts are alive
	http.HandleFunc(
		"/goloadbalance",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			encodedJson, err := json.Marshal(loadBalancer.GetServersStatus())
			if err != nil {
				log.Println(err)
			}
			w.Write(encodedJson)
		},
	)

	// Utility Function to add a server at runtime
	http.HandleFunc(
		"/goloadbalance/add_server",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusBadRequest)
			}

			requestBody, err := io.ReadAll(r.Body)
			var body AddServerRequest

			err = json.Unmarshal(requestBody, &body)

			if err != nil {
				log.Println("Error while parsing request body: ", err)
			}

			err = loadBalancer.AddServer(body.Host)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
			}

			log.Printf("Added %s to server pool", body.Host)
		},
	)

	// Utility Handler to check which hosts are alive
	http.HandleFunc(
		"/goloadbalance/remove_server",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusBadRequest)
			}

			body := RemoveServerRequest{}
			requestBody, err := io.ReadAll(r.Body)

			if err != nil {
				log.Println("Error while parsing request body: ", err.Error())
			}

			err = json.Unmarshal(requestBody, &body)

			if err != nil {
				log.Println("Error while parsing request body: ", err.Error())
			}

			parsedHostUrl, err := url.Parse(body.Host)

			if err != nil {
				log.Printf("Can't parse host url %s", body.Host)
			}

			server, err := loadBalancer.FindServer(parsedHostUrl)

			if err != nil {
				log.Printf(err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			log.Printf("Removing %s from the server pool", body.Host)
			loadBalancer.gracefullyShutdownServer(server)
		},
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type NoServerFoundErr struct{}

func (err *NoServerFoundErr) Error() string {
	return "No Server Found."
}

type ServerExists struct{}

func (err *ServerExists) Error() string {
	return "Server already exists."
}

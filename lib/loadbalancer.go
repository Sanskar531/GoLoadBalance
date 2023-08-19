package lib

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type LoadBalancer struct {
	Servers  []*Server
	balancer *Balancer
	cache    *Cache
}

func InitLoadBalancer(servers []*Server, balancer *Balancer) *LoadBalancer {
	return &LoadBalancer{
		Servers:  servers,
		balancer: balancer,
		cache:    InitCache(),
	}
}

func (loadBalancer *LoadBalancer) getServerToHandleRequest() *Server {
	return (*loadBalancer.balancer).GetServer(loadBalancer.Servers)
}

// This function is called as a go routine by the http module
// when serving a request
func (loadBalancer *LoadBalancer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	cachedMap := loadBalancer.cache.check(request)
	if cachedMap != nil {
		cachedResponse := (*cachedMap)["response"].(*http.Response)
		cachedBody := (*cachedMap)["body"].(*string)

		if cachedResponse != nil {
			log.Printf("Cache Hit: Found for route %s", request.URL.Path)

			// Reset the Headers properly before relaying the response back
			responseWriter.Header().Set("Content-Length", cachedResponse.Header.Get("Content-Length"))
			responseWriter.Header().Set("Content-Type", cachedResponse.Header.Get("Content-Type"))
			io.WriteString(responseWriter, *cachedBody)
			cachedResponse.Body.Close()
			return
		}
	}

	server := loadBalancer.getServerToHandleRequest()

	// Handle caching on a different thread so that we can return the response
	if res, body, err := server.HandleRequest(responseWriter, request); err == nil {
		go loadBalancer.cache.save(request, body, res)
	}
}

func (loadBalancer *LoadBalancer) GetServersStatus() map[string]bool {
	serverStatuses := make(map[string]bool)
	for _, server := range loadBalancer.Servers {
		serverStatuses[server.Url.Host] = server.Alive
	}
	return serverStatuses
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

	log.Fatal(http.ListenAndServe(":8080", nil))
}

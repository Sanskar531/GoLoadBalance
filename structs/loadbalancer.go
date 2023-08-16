package structs

import (
	"encoding/json"
	"log"
	"net/http"
)

type LoadBalancer struct {
	Servers  []*Server
	balancer *Balancer
}

func InitLoadBalancer(servers []*Server, balancer *Balancer) *LoadBalancer {
	return &LoadBalancer{
		Servers:  servers,
		balancer: balancer,
	}
}

func (loadBalancer *LoadBalancer) GetAliveServers() []*Server {
	var aliveServers []*Server

	// No need for a lock here because a single thread
	// will be updating whether the server is alive or not
	// as the healthcheck is a go routine
	for _, server := range loadBalancer.Servers {
		if server.Alive {
			aliveServers = append(aliveServers, server)
		}
	}

	return aliveServers
}

func (loadBalancer *LoadBalancer) getServerToHandleRequest() *Server {
	return (*loadBalancer.balancer).GetServer(loadBalancer.GetAliveServers())
}

// This function is called as a go routine by the http module
// when serving a request
func (loadBalancer *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := loadBalancer.getServerToHandleRequest()

	// Handle request on a different thread
	server.HandleRequest(w, r)
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

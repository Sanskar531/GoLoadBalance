package structs

import (
	"sync"
)

type LoadBalancer struct {
	Servers []Server
	mu sync.Mutex
}

func (loadBalancer *LoadBalancer) GetAliveBackends() []Server {
	var aliveBackends []Server;

	// A lock here because other threads might be pinging 
	// and changing status.
	loadBalancer.mu.Lock()
	for _, server := range loadBalancer.Servers {
		if server.Alive {
			aliveBackends = append(aliveBackends, server)
		}
	}
	loadBalancer.mu.Unlock()

	return aliveBackends;
}

func (loadbalancer *LoadBalancer) GetLoad() int {
	return len(loadbalancer.Servers)
}

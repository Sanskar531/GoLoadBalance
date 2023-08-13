package structs

import (
	"sync"
)

type LoadBalancer struct {
	Servers  []Server
	mu       sync.Mutex
	balancer Balancer
}

func InitLoadBalancer(servers []Server, balancer Balancer) LoadBalancer {
	return LoadBalancer{
		Servers:  servers,
		balancer: balancer,
	}
}

func (loadBalancer *LoadBalancer) GetAliveBackends() []Server {
	var aliveBackends []Server

	// A lock here because other threads might be pinging
	// and changing status.
	loadBalancer.mu.Lock()
	for _, server := range loadBalancer.Servers {
		if server.Alive {
			aliveBackends = append(aliveBackends, server)
		}
	}
	loadBalancer.mu.Unlock()

	return aliveBackends
}

func (loadBalancer *LoadBalancer) getBackendToServe() Server {
	return loadBalancer.balancer.GetServer(loadBalancer.GetAliveBackends())
}

func (loadBalancer *LoadBalancer) Balance() {
	for {
		loadBalancer.getBackendToServe()
	}
}

package structs

import (
	"container/ring"
)

type Balancer interface {
	GetServer(servers []*Server) *Server
}

type RoundRobin struct {
	ringBuffer *ring.Ring
}

func InitRoundRobin(serverPoolSize int) Balancer {
	// Initialize ring buffer so we can cycle through
	// the servers in a round robin without needing to
	// keep a index.
	ringBuffer := ring.New(serverPoolSize)

	// Initialize the values of the ring buffer
	// to it's index.
	for i := 0; i < serverPoolSize; i++ {
		ringBuffer.Value = i
		ringBuffer = ringBuffer.Next()
	}

	return &RoundRobin{
		ringBuffer: ringBuffer,
	}
}

func (balancer *RoundRobin) GetServer(servers []*Server) *Server {
	// Round robin works by distributing the load equally to incoming requests
	// hence here we just get the value of the current value in the index of the
	// ring buffer select the server we want to supply the load to and then return
	// it

	currentIndex := balancer.ringBuffer.Value
	balancer.ringBuffer = balancer.ringBuffer.Next()
	return servers[currentIndex.(int)]
}

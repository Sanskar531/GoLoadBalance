package lib

import (
	"container/ring"
	"sync"
)

type Balancer interface {
	GetServer(servers []*Server) *Server
	ServerDead()
	ServerAdd()
}

type RoundRobin struct {
	ringBuffer *ring.Ring
	mutex      *sync.Mutex
}

func InitRoundRobin(serverPoolSize int) Balancer {
	// Initialize ring buffer so we can cycle through
	// the servers in a round robin without needing to
	// keep a index.
	ringBuffer := initRingBuffer(serverPoolSize)

	return &RoundRobin{
		ringBuffer: ringBuffer,
		mutex:      &sync.Mutex{},
	}
}

func initRingBuffer(size int) *ring.Ring {
	ringBuffer := ring.New(size)

	// Initialize the values of the ring buffer
	// to it's index.
	for i := 0; i < size; i++ {
		ringBuffer.Value = i
		ringBuffer = ringBuffer.Next()
	}

	return ringBuffer
}

func (balancer *RoundRobin) ServerDead() {
	balancer.mutex.Lock()
	defer balancer.mutex.Unlock()

	// Once server died so ring buffer needs to go down by one.
	balancer.ringBuffer = initRingBuffer(balancer.ringBuffer.Len() - 1)
}

func (balancer *RoundRobin) ServerAdd() {
	balancer.mutex.Lock()
	defer balancer.mutex.Unlock()

	// Once server died so ring buffer needs to go down by one.
	balancer.ringBuffer = initRingBuffer(balancer.ringBuffer.Len() + 1)
}

func (balancer *RoundRobin) GetServer(servers []*Server) *Server {
	// Round robin works by distributing the load equally to incoming requests
	// hence here we just get the value of the current value in the index of the
	// ring buffer select the server we want to supply the load to and then return
	// it

	// Accessing the ring buffer is not thread safe so we need a mutex
	balancer.mutex.Lock()
	defer balancer.mutex.Unlock()

	currentIndex := balancer.ringBuffer.Value
	balancer.ringBuffer = balancer.ringBuffer.Next()
	return servers[currentIndex.(int)]
}

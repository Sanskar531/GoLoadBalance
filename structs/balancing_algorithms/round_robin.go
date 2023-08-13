package balancingalgorithms

import (
	"container/ring"
	"github.com/sanskar531/goloadbalance/structs"
)

type RoundRobin struct {
	ringBuffer *ring.Ring
}

func InitRoundRobin(serverPoolSize int) structs.Balancer {
	ringBuffer := ring.New(serverPoolSize)

	for i := 0; i < serverPoolSize; i++ {
		ringBuffer.Value = i
		ringBuffer = ringBuffer.Next()
	}

	return &RoundRobin{
		ringBuffer: ringBuffer,
	}
}

func (balancer *RoundRobin) GetServer(servers []structs.Server) structs.Server {
	currentIndex := balancer.ringBuffer.Value
	balancer.ringBuffer = balancer.ringBuffer.Next()
	return servers[currentIndex.(int)]
}

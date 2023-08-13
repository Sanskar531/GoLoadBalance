package balancingalgorithms

import (
	"fmt"
	"github.com/sanskar531/goloadbalance/structs"
)

type RoundRobin struct {
	currentIndex int
}

func InitRoundRobin () structs.Balancer{
	roundRobin := new(RoundRobin);
	return roundRobin;
}

func (balancer *RoundRobin) GetServer(servers []structs.Server) structs.Server {
	fmt.Println("HEREE!")
	return servers[0]
}

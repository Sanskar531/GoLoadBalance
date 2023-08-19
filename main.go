package main

import (
	"github.com/sanskar531/goloadbalance/lib"
)

func main() {
	config := parseCommandLineArgs()

	var servers []*lib.Server

	for _, parsedUrl := range config.ServerUrls {
		servers = append(servers, lib.InitServer(
			parsedUrl,
		))
	}

	balancer := lib.InitRoundRobin(len(servers))

	loadbalancer := lib.InitLoadBalancer(
		servers,
		&balancer,
	)

	loadbalancer.Balance()
}

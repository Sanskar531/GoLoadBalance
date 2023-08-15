package main

import (
	"github.com/sanskar531/goloadbalance/structs"
)

func main() {
	parsedUrls, _ := parseCommandLineArgs()

	var servers []*structs.Server

	for _, parsedUrl := range parsedUrls {
		servers = append(servers, structs.InitServer(
			parsedUrl,
		))
	}

	balancer := structs.InitRoundRobin(len(servers))

	loadbalancer := structs.InitLoadBalancer(
		servers,
		balancer,
	)

	loadbalancer.Balance()
}

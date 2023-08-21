package main

import (
	"github.com/sanskar531/goloadbalance/lib"
	"github.com/sanskar531/goloadbalance/configuration"
)

func main() {
	mainConfig := configuration.ParseCommandLineArgs()

	var servers []*lib.Server

	for _, parsedUrl := range mainConfig.ServerUrls {
		servers = append(servers, lib.InitServer(
			parsedUrl,
			mainConfig,
		))
	}

	balancer := lib.InitRoundRobin(len(servers))

	loadbalancer := lib.InitLoadBalancer(
		servers,
		&balancer,
		mainConfig.CacheEnabled,
		mainConfig.CacheTimeoutInSeconds,
	)

	loadbalancer.Balance()
}

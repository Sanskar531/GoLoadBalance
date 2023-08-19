package main

import (
	"github.com/sanskar531/goloadbalance/lib"
)

var config *Config

func main() {
	config = parseCommandLineArgs()

	var servers []*lib.Server

	for _, parsedUrl := range config.ServerUrls {
		servers = append(servers, lib.InitServer(
			parsedUrl,
			config.HealthCheckFrequencyInSeconds,
		))
	}

	balancer := lib.InitRoundRobin(len(servers))

	loadbalancer := lib.InitLoadBalancer(
		servers,
		&balancer,
		config.CacheEnabled,
		config.CacheTimeoutInSeconds,
	)

	loadbalancer.Balance()
}

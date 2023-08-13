package main

import (
	"flag"
	"fmt"
	"github.com/sanskar531/goloadbalance/structs"
	"github.com/sanskar531/goloadbalance/structs/balancing_algorithms"
	"net/url"
)

var help = flag.Bool("help", false, "Show help")

func parseCommandLineArgs() {
	flag.Parse()

	if *help {
		fmt.Print(`
			GOPrintlnPrintlnPrintln
		`)
		return
	}
}

func main() {
	parseCommandLineArgs()

	serverUrl, err := url.Parse("www.google.com")
	if err != nil {
		return
	}

	balancer := balancingalgorithms.InitRoundRobin()

	loadbalancer := structs.InitLoadBalancer(
		[]structs.Server{
			structs.InitServer(
				serverUrl,
			),
		},
		balancer,
	)
	go loadbalancer.Balance()
	fmt.Scanln()
}

package main

import (
	"flag"
	"fmt"
	"github.com/sanskar531/goloadbalance/structs"
	"github.com/sanskar531/goloadbalance/structs/balancing_algorithms"
	"net/url"
)


func parseCommandLineArgs() {
	help := flag.Bool("help", false, "Show help")
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

	serverUrl, err := url.Parse("http://localhost:8000")
	if err != nil {
		return
	}

	servers := []structs.Server{
			structs.InitServer(
				serverUrl,
			),
	};

	balancer := balancingalgorithms.InitRoundRobin(len(servers))

	loadbalancer := structs.InitLoadBalancer(
		servers,
		balancer,
	)

	go loadbalancer.Balance()
	fmt.Scanln()
}

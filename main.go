package main

import (
	"flag"
	"fmt"
	"net/url"
	"github.com/sanskar531/goloadbalance/structs"
)


var help = flag.Bool("help", false, "Show help")

func parseCommandLineArgs() {
	flag.Parse();

	if *help {
		fmt.Print(`
			GOPrintlnPrintlnPrintln
		`)
		return
	}
}

func main() {
	parseCommandLineArgs();

	serverUrl, err := url.Parse("www.google.com")
	if err != nil {
		return
	}
	loadbalancer := structs.LoadBalancer{
		Servers: []structs.Server{
			structs.InitServer(
				serverUrl,
			),
		},
	}

	fmt.Println(loadbalancer.GetLoad())
	fmt.Println(loadbalancer.GetAliveBackends())
	fmt.Scanln();
}

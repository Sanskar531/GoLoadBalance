package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/sanskar531/goloadbalance/structs"
	"github.com/sanskar531/goloadbalance/structs/balancing_algorithms"
)


func parseCommandLineArgs() []*url.URL {
	servers := flag.String("servers", "", "Show servers");
	flag.Parse();

	if *servers == "" {
		log.Fatal("Please enter at least one server.")
	}

	serverUrls := strings.Split(*servers, ",")
	var parsedUrls []*url.URL;

	for _, serverUrl := range serverUrls {
		parsedUrl, err := url.Parse(serverUrl);
		if err != nil {
			log.Fatal("The following url is invalid: ", parsedUrl)
		}

		parsedUrls = append(parsedUrls, parsedUrl)
	}

	return parsedUrls
}

func main() {
	parsedUrls := parseCommandLineArgs();

	var servers []structs.Server;

	for _, parsedUrl := range parsedUrls {
		servers = append(servers, structs.InitServer(
			parsedUrl,
		))
	}

	balancer := balancingalgorithms.InitRoundRobin(len(servers))

	loadbalancer := structs.InitLoadBalancer(
		servers,
		balancer,
	)

	go loadbalancer.Balance()
	fmt.Scanln()
}

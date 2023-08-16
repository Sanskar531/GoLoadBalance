package main

import (
	"flag"
	"log"
	"net/url"
	"strings"
)

type Config struct {
	ServerUrls []*url.URL
	Algorithm  string
}

func parseServers(servers *string) []*url.URL {
	if *servers == "" {
		log.Fatal("Please enter at least one server.")
	}

	serverUrls := strings.Split(*servers, ",")
	var parsedUrls []*url.URL

	for _, serverUrl := range serverUrls {
		parsedUrl, err := url.Parse(serverUrl)
		if err != nil {
			log.Fatal("The following url is invalid: ", parsedUrl)
		}

		parsedUrls = append(parsedUrls, parsedUrl)
	}

	return parsedUrls
}

func parseBalancingAlgorithm(algorithm *string) string {

	if *algorithm == "round_robin" {
		return "round_robin"
	}

	return "round_robin"
}

func parseCommandLineArgs() *Config {
	servers := flag.String("servers", "", "Server urls. Usage: --servers=http://localhost:3000")
	algorithm := flag.String("algorithm", "", "Load Balancing Algorithms: round_robin. Usage: --algorithm=round_robin")
	configFilePath := flag.String("config", "", "Config file for the load balancer. Usage --config=./example.yaml")

	flag.Parse()

	if (*configFilePath) != "" {
		return parseYAML(*configFilePath)
	}

	return &Config{
		ServerUrls: parseServers(servers),
		Algorithm:  parseBalancingAlgorithm(algorithm),
	}
}

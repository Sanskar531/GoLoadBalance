package main

import (
	"flag"
	"log"
	"net/url"
	"strings"

)

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

func parseCommandLineArgs() ([]*url.URL, string) {
	servers := flag.String("servers", "", "Show servers")
	algorithm := flag.String("algorithm", "", "Show Algorithms")

	flag.Parse()
	return parseServers(servers), parseBalancingAlgorithm(algorithm)
}

package main

import (
	"flag"
	"log"
	"net/url"
	"strings"
)

type Config struct {
	ServerUrls                    []*url.URL
	Algorithm                     string
	CacheEnabled                  bool
	CacheTimeoutInSeconds         int
	HealthCheckFrequencyInSeconds int
	HealthCheckMaxRetries         int
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
	configFilePath := flag.String("config", "", "Config file for the load balancer. Usage: --config=./example.yaml")
	cachingEnabled := flag.Bool("cache-enabled", false, "Enable Caching in the load balancer. Usage: --cache-enabled")
	cacheTimout := flag.Int("cache-timeout-in-seconds", 0, "Keep cached value alive for x seconds in the load balancer. Usage: --cache-timeout=10")
	healthCheckFrequency := flag.Int("health-check-frequency-in-seconds", 0, "The amount of time in between a ping to the server. Usage: --health-check-frequency-in-seconds=10")
	healthCheckMaxRetry := flag.Int("health-check-max-retries", 0, "The amount of times pings will occur after initially when a server doesn't respond. Usage: --health-check-max-retries=10")

	flag.Parse()

	if (*configFilePath) != "" {
		return parseYAML(*configFilePath)
	}

	return &Config{
		ServerUrls:                    parseServers(servers),
		Algorithm:                     parseBalancingAlgorithm(algorithm),
		CacheEnabled:                  *cachingEnabled,
		CacheTimeoutInSeconds:         *cacheTimout,
		HealthCheckFrequencyInSeconds: *healthCheckFrequency,
		HealthCheckMaxRetries:         *healthCheckMaxRetry,
	}
}

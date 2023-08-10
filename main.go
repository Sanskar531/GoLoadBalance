package main

import "github.com/sanskar531/goloadbalance/structs"

func main() {
	loadbalancer := structs.LoadBalancer{
		Workers: []structs.Worker{
			{
				CurrentNumberOfRequests: 12,
			},
		},
	}
}

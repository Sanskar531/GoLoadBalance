package structs

type LoadBalancer struct {
	Workers []Worker
}

func (loadbalancer LoadBalancer) getLoad() string {
	return "Lol"
}

package structs

type Balancer interface {
	GetServer(servers []Server) Server
}

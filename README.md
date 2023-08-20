# GoLoadBalance

This is a personal project of mine to try and understand how requests are handled by a load balancer and also get more familiar with GO, go routines and it's STD library. A loadbalancer is one of the first place where your request hits so, having a deep understanding of such infrastructure will help a lot. Additionally, implementing one from scratch means that we will have multiple requests going to different servers concurently whilst also doing health check to make sure that our server is fine which will help in honing my concurrent programming skills. 

The main goal of this load balancer is to first be initialized and be open to receiving requests on a specific port. Once the port is specified, it will then receive the request and then redirect the request appropriately to the selected servers using a load balancing algorithm such as round robin.

---

## Things implemented:
- [x] Can pass servers urls/(IPS+port) via cli.
- [x] Can pass balancing algorithms via cli.
- [x] Health checks to ping server whether it's still alive or not.
- [x] Actual loading balancing using `round_robin` and `Circular/Ring buffers`. 
- [x] Support passing `.yaml` config files
- [x] Basic Caching by hashing parts of the request header
- [x] Max Retries before a server is removed from the server pool.

## Things left to implement:
- [ ] TLS support
- [ ] Permanent redirects/rewrites like nginx
- [ ] Support multiple load balancing algorithm
- [ ] Feature for server notifying the balancer that it is in maintainance
- [ ] Support Webhooks for when a server dies
- [ ] Add/Remove server at runtime
- [ ] Support least connections loadbalancing algorithm.

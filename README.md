# GoLoadBalance

This is a personal project of mine to try and understand how requests are handled by a load balancer and also get more familiar with GO, go routines and it's STD library. A loadbalancer is one of the first place where your request hits so, having a deep understanding of such infrastructure will help a lot. Additionally, implementing one from scratch means that we will have multiple requests going to different servers concurently whilst also doing health check to make sure that our server is fine which will help in honing my concurrent programming skills. 

The main goal of this load balancer is to first be initialized and be open to receiving requests on a specific port. Once the port is specified, it will then receive the request and then redirect the request appropriately to the selected servers using a load balancing algorithm such as round robin.

---

## Things implemented:
- [x] Can pass servers urls/(IPS+port) via cli.
- [x] Can pass balancing algorithms via cli.
- [x] Health checks to ping server whether it's still alive or not.
- [x] Actual loading balancing using `round_robin` and `circular buffers`. 
- [x] Support passing `.yaml` config files
- [x] Caching

## Things left to implement:
- [ ] TLS support
- [ ] Permanent redirects/rewrites like nginx
- [ ] Support multiple load balancing algorithm

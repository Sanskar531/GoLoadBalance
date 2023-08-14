# GoLoadBalance

A personal project of mine to try and understand how requests are handled by a load balancer and also get more familiar with go and it's std library. A Loadbalancer is one of the first place where your request hits. 

The main goal of this load balancer is to first be initialized and be open to receiving requests on a specific port. Once the port is specified, it will then receive the request and then redirect the request appropriately to the selected servers using.

---

## Things implemented:
    - [x] Can pass servers urls/(IPS+port) via cli.
    - [x] Can pass balancing algorithms via cli.
    - [x] Health checks to ping server whether it's still alive or not.

## Things left to implement:
    - [] Caching
    - [] TLS support
    - [] Permanent redirects/rewrites like nginx
    - [] Support passing `.toml` config files

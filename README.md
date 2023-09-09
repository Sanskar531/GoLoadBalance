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
- [x] Add/Remove server at runtime
- [x] Blacklisting IPs
- [x] Support Webhooks for when a server dies

## Things left to implement:
- [ ] TLS support.
- [ ] Permanent redirects/rewrites like nginx.
- [ ] Support multiple load balancing algorithms.
- [ ] Use channels to optimize to avoid mutexes.
- [ ] Use go routine pool so we don't need to spawn so many goroutines for many concurrent request.

---

## License

Copyright 2023 Sanskar Gauchan

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

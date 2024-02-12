# FastHTTP Load Balancing Client

Flexible replacement of fasthttp's [LBClient](https://github.com/valyala/fasthttp/blob/master/lbclient.go).

Our experience: we have ~20 instances of the service on different servers.
We use [LBClient](https://github.com/valyala/fasthttp/blob/master/lbclient.go#L27) together with
[HostClient](https://github.com/valyala/fasthttp/blob/master/client.go#L607) to balance requests among them.
Average RPS is 3-5k requests.

Sometime one of these clients become unavailable, mostly due to deployment process.
HostClient instance of this server collects zero pending requests and maximum 300 penalty counter.

In result this [condition](https://github.com/valyala/fasthttp/blob/master/lbclient.go#L142) decides that
this server is least loaded and prioritize it. Almost 100% of our traffic drains into the trash.

This package was developed to solve that problem and also provides few handy features, see next section.

## Features

* Health check function replaced to interface with corresponding function.
* Added Penalty field to handle penalty duration.
* Added Balancer field where you can specify your own balance algorithm, see [existing](https://github.com/koykov/fhlbclient/tree/master/balancer).
* Added RequestHooker field where you can specify implementation of PreRequest() and PostRequest() methods.

package fhlbclient

import (
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type HealthCheckFn func(req *fasthttp.Request, resp *fasthttp.Response, err error) bool

type LBClient struct {
	Clients []fasthttp.BalancingClient
	HealthCheck HealthCheckFn
	Timeout time.Duration
	Balancer Balancer
	once sync.Once

	cln []innerClient
}

type innerClient struct {
	bc fasthttp.BalancingClient
	hc HealthCheckFn
	pen int32
	tot int64
}

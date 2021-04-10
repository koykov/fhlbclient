package fhlbclient

import (
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	DefaultTimeout = time.Second
	DefaultPenalty = time.Second * 3

	maxPenalty = 300
)

type HealthCheckFn func(req *fasthttp.Request, resp *fasthttp.Response, err error) bool

type LBClient struct {
	Clients     []fasthttp.BalancingClient
	HealthCheck HealthCheckFn
	Timeout     time.Duration
	Penalty     time.Duration
	Balancer    Balancer
	once        sync.Once

	cln []innerClient
}

func (c *LBClient) init() {
	c.cln = make([]innerClient, 0, len(c.Clients))
	for _, bc := range c.Clients {
		pd := c.Penalty
		if pd <= 0 {
			pd = DefaultPenalty
		}
		c.cln = append(c.cln, innerClient{
			bc: bc,
			hc: c.HealthCheck,
			pd: pd,
		})
	}
}

func (c *LBClient) DoDeadline(req *fasthttp.Request, resp *fasthttp.Response, deadline time.Time) error {
	return c.get().DoDeadline(req, resp, deadline)
}

func (c *LBClient) DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	return c.get().DoDeadline(req, resp, deadline)
}

func (c *LBClient) Do(req *fasthttp.Request, resp *fasthttp.Response) error {
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return c.DoTimeout(req, resp, timeout)
}

func (c *LBClient) get() *innerClient {
	c.once.Do(c.init)
	return nil
}

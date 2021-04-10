package fhlbclient

import (
	"errors"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	DefaultTimeout = time.Second
	DefaultPenalty = time.Second * 3
)

type HealthCheckFn func(req *fasthttp.Request, resp *fasthttp.Response, err error) bool

type LBClient struct {
	Clients     []fasthttp.BalancingClient
	HealthCheck HealthCheckFn
	Timeout     time.Duration
	Penalty     time.Duration
	Balancer    Balancer
	once        sync.Once

	cln []PenalizingClient
}

var ErrNoAliveClients = errors.New("no alive clients available")

func (c *LBClient) init() {
	pd := c.Penalty
	if pd <= 0 {
		pd = DefaultPenalty
	}
	c.cln = make([]PenalizingClient, 0, len(c.Clients))
	for _, bc := range c.Clients {
		c.cln = append(c.cln, PenalizingClient{
			bc: bc,
			hc: c.HealthCheck,
			pd: pd,
		})
	}
}

func (c *LBClient) DoDeadline(req *fasthttp.Request, resp *fasthttp.Response, deadline time.Time) error {
	if pc := c.get(); pc != nil {
		return pc.DoDeadline(req, resp, deadline)
	}
	return ErrNoAliveClients
}

func (c *LBClient) DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration) error {
	if pc := c.get(); pc != nil {
		deadline := time.Now().Add(timeout)
		return pc.DoDeadline(req, resp, deadline)
	}
	return ErrNoAliveClients
}

func (c *LBClient) Do(req *fasthttp.Request, resp *fasthttp.Response) error {
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return c.DoTimeout(req, resp, timeout)
}

func (c *LBClient) get() *PenalizingClient {
	c.once.Do(c.init)
	if len(c.cln) == 0 {
		return nil
	}
	return c.Balancer.Evaluate(c.cln)
}

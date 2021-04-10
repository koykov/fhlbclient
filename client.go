package fhlbclient

import (
	"errors"
	"sync"
	"sync/atomic"
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
	pc := c.get()
	if pc == nil {
		return ErrNoAliveClients
	}
	return pc.DoDeadline(req, resp, deadline)
}

func (c *LBClient) DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration) error {
	pc := c.get()
	if pc == nil {
		return ErrNoAliveClients
	}
	deadline := time.Now().Add(timeout)
	return pc.DoDeadline(req, resp, deadline)
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

	var (
		minC *PenalizingClient
		off  int
	)
	for i := 0; i < len(c.cln); i++ {
		if atomic.LoadInt32(&c.cln[i].pen) == 0 {
			minC = &c.cln[i]
			off = i + 1
			break
		}
	}
	if minC == nil {
		return nil
	}

	if off < len(c.cln) {
		minN, minT := minC.RequestStats()
		for i := off; i < len(c.cln); i++ {
			pc := &c.cln[i]
			if atomic.LoadInt32(&pc.pen) > 0 {
				continue
			}
			n, t := pc.RequestStats()
			if n < minN || (n == minN && t < minT) {
				minC = pc
				minN = n
				minT = t
			}
		}
	}

	return minC
}

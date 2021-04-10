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

	cln []innerClient
}

var ErrNoAliveClients = errors.New("no alive clients available")

func (c *LBClient) init() {
	pd := c.Penalty
	if pd <= 0 {
		pd = DefaultPenalty
	}
	c.cln = make([]innerClient, 0, len(c.Clients))
	for _, bc := range c.Clients {
		c.cln = append(c.cln, innerClient{
			bc: bc,
			hc: c.HealthCheck,
			pd: pd,
		})
	}
}

func (c *LBClient) DoDeadline(req *fasthttp.Request, resp *fasthttp.Response, deadline time.Time) error {
	ic := c.get()
	if ic == nil {
		return ErrNoAliveClients
	}
	return ic.DoDeadline(req, resp, deadline)
}

func (c *LBClient) DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration) error {
	ic := c.get()
	if ic == nil {
		return ErrNoAliveClients
	}
	deadline := time.Now().Add(timeout)
	return ic.DoDeadline(req, resp, deadline)
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

	var (
		minC *innerClient
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
		minN := minC.PendingRequests()
		minT := atomic.LoadUint64(&minC.tot)
		for i := off; i < len(c.cln); i++ {
			ic := &c.cln[i]
			if atomic.LoadInt32(&ic.pen) > 0 {
				continue
			}
			n := ic.PendingRequests()
			t := atomic.LoadUint64(&ic.tot)
			if n < minN || (n == minN && t < minT) {
				minC = ic
				minN = n
				minT = t
			}
		}
	}

	return minC
}

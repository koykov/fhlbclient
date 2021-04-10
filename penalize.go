package fhlbclient

import (
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
)

type PenalizingClient struct {
	bc  fasthttp.BalancingClient
	hc  HealthCheckFn
	pd  time.Duration
	pen int32
	tot uint64
}

func (c *PenalizingClient) DoDeadline(req *fasthttp.Request, resp *fasthttp.Response, deadline time.Time) error {
	err := c.bc.DoDeadline(req, resp, deadline)
	if !c.isHealthy(req, resp, err) && c.incPenalty() {
		time.AfterFunc(c.pd, c.decPenalty)
	} else {
		atomic.AddUint64(&c.tot, 1)
	}
	return err
}

func (c *PenalizingClient) RequestStats() (uint64, uint64) {
	return uint64(c.bc.PendingRequests() + int(atomic.LoadInt32(&c.pen))), atomic.LoadUint64(&c.tot)
}

func (c *PenalizingClient) isHealthy(req *fasthttp.Request, resp *fasthttp.Response, err error) bool {
	if c.hc == nil {
		return err == nil
	}
	return c.hc(req, resp, err)
}

func (c *PenalizingClient) incPenalty() bool {
	m := atomic.AddInt32(&c.pen, 1)
	if m > maxPenalty {
		c.decPenalty()
		return false
	}
	return true
}

func (c *PenalizingClient) decPenalty() {
	atomic.AddInt32(&c.pen, -1)
}

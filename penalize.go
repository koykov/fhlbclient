package fhlbclient

import (
	"math"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
)

const maxPenalty = 300

// A wrapper around BalancingClient with penalty support.
type PenalizingClient struct {
	// BalancingClient instance.
	bc fasthttp.BalancingClient
	// Health checker helper.
	hc HealthChecker
	// Penalty duration.
	pd time.Duration
	// Penalty counter.
	// pen > 0 means that client is under penalty.
	pen uint32
	// Total requests counter.
	tot uint64
}

// Execute request with given deadline.
func (c *PenalizingClient) DoDeadline(req *fasthttp.Request, resp *fasthttp.Response, deadline time.Time) error {
	// Execute request.
	err := c.bc.DoDeadline(req, resp, deadline)
	// Check health state and increase penalty counter.
	if !c.isHealthy(req, resp, err) && c.incPenalty() {
		// Register postponed func to decrease the counter.
		time.AfterFunc(c.pd, c.decPenalty)
	} else {
		// Increase total counter.
		atomic.AddUint64(&c.tot, 1)
	}
	return err
}

// Get two requests metrics: pending requests and total requests counts.
//
// Pending requests value includes penalty counter value.
func (c *PenalizingClient) RequestStats() (uint64, uint64) {
	return uint64(c.bc.PendingRequests() + int(atomic.LoadUint32(&c.pen))), atomic.LoadUint64(&c.tot)
}

// Check if client is under penalty.
func (c *PenalizingClient) UnderPenalty() bool {
	return atomic.LoadUint32(&c.pen) > 0
}

// Get inner fasthttp's balancing client instance.
func (c *PenalizingClient) Instance() fasthttp.BalancingClient {
	return c.bc
}

// Check if client has good health.
func (c *PenalizingClient) isHealthy(req *fasthttp.Request, resp *fasthttp.Response, err error) bool {
	if c.hc == nil {
		return err == nil
	}
	return c.hc.Check(req, resp, err)
}

// Increase penalty counter.
func (c *PenalizingClient) incPenalty() bool {
	if m := atomic.AddUint32(&c.pen, 1); m > maxPenalty {
		c.decPenalty()
		return false
	}
	return true
}

// Decrease penalty counter.
func (c *PenalizingClient) decPenalty() {
	atomic.AddUint32(&c.pen, math.MaxUint32)
}

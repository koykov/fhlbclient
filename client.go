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

// Load balancing client.
//
// see https://github.com/valyala/fasthttp/blob/master/lbclient.go for details and comparison.
type LBClient struct {
	// Array of clients to balance.
	Clients []fasthttp.BalancingClient
	// Health check helper.
	HealthCheck HealthChecker
	// Timeout duration to execute request.
	// Will used DefaultTimeout if empty.
	Timeout time.Duration
	// Penalty duration to ban the client.
	Penalty time.Duration
	// Balancer helper.
	Balancer Balancer
	// Array of wrappers around Clients.
	// Note, that Clients used only for init step and copies into cln.
	cln []PenalizingClient

	once sync.Once
}

var ErrNoAliveClients = errors.New("no alive clients available")

// Init the load balancing client.
func (c *LBClient) init() {
	// Get penalty duration.
	pd := c.Penalty
	if pd <= 0 {
		pd = DefaultPenalty
	}
	// Make new PenalizingClient for each provided BalancingClient.
	c.cln = make([]PenalizingClient, 0, len(c.Clients))
	for _, bc := range c.Clients {
		c.cln = append(c.cln, PenalizingClient{
			bc: bc,
			hc: c.HealthCheck,
			pd: pd,
		})
	}
}

// Execute request with given deadline.
func (c *LBClient) DoDeadline(req *fasthttp.Request, resp *fasthttp.Response, deadline time.Time) error {
	if pc := c.get(); pc != nil {
		return pc.DoDeadline(req, resp, deadline)
	}
	// No available clients found (all of them under penalty).
	return ErrNoAliveClients
}

// Execute request with given timeout.
func (c *LBClient) DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration) error {
	if pc := c.get(); pc != nil {
		deadline := time.Now().Add(timeout)
		return pc.DoDeadline(req, resp, deadline)
	}
	// No available clients found (all of them under penalty).
	return ErrNoAliveClients
}

// Execute request with internal timeout.
func (c *LBClient) Do(req *fasthttp.Request, resp *fasthttp.Response) error {
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return c.DoTimeout(req, resp, timeout)
}

// Get least loaded client.
func (c *LBClient) get() *PenalizingClient {
	// Run init() once.
	c.once.Do(c.init)
	if len(c.cln) == 0 {
		return nil
	}
	// Use balancer helper to get best candidate.
	return c.Balancer.Evaluate(c.cln)
}

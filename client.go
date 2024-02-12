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

// LBClient implements load balancing client.
//
// See https://github.com/valyala/fasthttp/blob/master/lbclient.go for details and comparison.
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
	// Request hooks helper.
	RequestHooker RequestHooker
	// Metrics writer handler.
	// Available only in balancing methods: DoDeadlineWB, DoTimeoutWB and DoWB.
	MetricsWriter MetricsWriter
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
	// Check balancer helper.
	if c.Balancer == nil {
		c.Balancer = DummyBalancer{}
	}
	// Check request hooks helper.
	if c.RequestHooker == nil {
		c.RequestHooker = DummyRequestHooks{}
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
	// Check metrics writer.
	if c.MetricsWriter == nil {
		c.MetricsWriter = DummyMetricsWriter{}
	}
}

// DoDeadline executes request with given deadline.
func (c *LBClient) DoDeadline(req *fasthttp.Request, resp *fasthttp.Response, deadline time.Time) error {
	return c.DoDeadlineWB(req, resp, deadline, c.Balancer)
}

// DoTimeout executes request with given timeout.
func (c *LBClient) DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration) error {
	return c.DoTimeoutWB(req, resp, timeout, c.Balancer)
}

// Do executes request with internal timeout.
func (c *LBClient) Do(req *fasthttp.Request, resp *fasthttp.Response) error {
	return c.DoWB(req, resp, c.Balancer)
}

// Get least loaded client.
func (c *LBClient) get(b Balancer) *PenalizingClient {
	// Run init() once.
	c.once.Do(c.init)
	if len(c.cln) == 0 || b == nil {
		return nil
	}
	// Use balancer helper to get best candidate.
	return b.Evaluate(c.cln)
}

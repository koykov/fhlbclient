package fhlbclient

// Collection of Do*() methods with explicit balancer.
// Need for cases then balancer contains some data that shouldn't be shared among goroutines.

import (
	"time"

	"github.com/koykov/fastconv"
	"github.com/valyala/fasthttp"
)

// DoDeadlineWB executes request with given deadline and balancer.
func (c *LBClient) DoDeadlineWB(req *fasthttp.Request, resp *fasthttp.Response, deadline time.Time, b Balancer) error {
	if pc := c.get(b); pc != nil {
		c.RequestHooker.PreRequest(req, resp, pc)
		err := pc.DoDeadline(req, resp, deadline)
		c.RequestHooker.PostRequest(req, resp, pc, err)
		c.MetricsWriter.HostStatus(fastconv.B2S(req.Host()), resp.StatusCode())
		return err
	}
	// No available clients found (all of them under penalty).
	c.MetricsWriter.HostThrottle(fastconv.B2S(req.Host()))
	return ErrNoAliveClients
}

// DoTimeoutWB executes request with given timeout and balancer.
func (c *LBClient) DoTimeoutWB(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration, b Balancer) error {
	if pc := c.get(b); pc != nil {
		deadline := time.Now().Add(timeout)
		c.RequestHooker.PreRequest(req, resp, pc)
		err := pc.DoDeadline(req, resp, deadline)
		c.RequestHooker.PostRequest(req, resp, pc, err)
		c.MetricsWriter.HostStatus(fastconv.B2S(req.Host()), resp.StatusCode())
		return err
	}
	// No available clients found (all of them under penalty).
	c.MetricsWriter.HostThrottle(fastconv.B2S(req.Host()))
	return ErrNoAliveClients
}

// DoWB executes request with internal timeout and balancer.
func (c *LBClient) DoWB(req *fasthttp.Request, resp *fasthttp.Response, b Balancer) error {
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return c.DoTimeoutWB(req, resp, timeout, b)
}

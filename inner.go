package fhlbclient

import (
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
)

type innerClient struct {
	bc  fasthttp.BalancingClient
	hc  HealthCheckFn
	pd  time.Duration
	pen int32
	tot uint64
}

func (ic *innerClient) DoDeadline(req *fasthttp.Request, resp *fasthttp.Response, deadline time.Time) error {
	err := ic.bc.DoDeadline(req, resp, deadline)
	if !ic.isHealthy(req, resp, err) && ic.incPenalty() {
		time.AfterFunc(ic.pd, ic.decPenalty)
	} else {
		atomic.AddUint64(&ic.tot, 1)
	}
	return err
}

func (ic *innerClient) PendingRequests() int {
	n := ic.bc.PendingRequests()
	m := atomic.LoadInt32(&ic.pen)
	return n + int(m)
}

func (ic *innerClient) isHealthy(req *fasthttp.Request, resp *fasthttp.Response, err error) bool {
	if ic.hc == nil {
		return err == nil
	}
	return ic.hc(req, resp, err)
}

func (ic *innerClient) incPenalty() bool {
	m := atomic.AddInt32(&ic.pen, 1)
	if m > maxPenalty {
		ic.decPenalty()
		return false
	}
	return true
}

func (ic *innerClient) decPenalty() {
	atomic.AddInt32(&ic.pen, -1)
}

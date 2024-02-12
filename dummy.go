package fhlbclient

import "github.com/valyala/fasthttp"

type DummyBalancer struct{}

func (DummyBalancer) Evaluate(_ []PenalizingClient) *PenalizingClient { return nil }

type DummyRequestHooks struct{}

func (DummyRequestHooks) PreRequest(_ *fasthttp.Request, _ *fasthttp.Response, _ *PenalizingClient) {}
func (DummyRequestHooks) PostRequest(_ *fasthttp.Request, _ *fasthttp.Response, _ *PenalizingClient, _ error) {
}

type DummyMetricsWriter struct{}

func (DummyMetricsWriter) HostStatus(_ string, _ int) {}
func (DummyMetricsWriter) HostThrottle(_ string)      {}

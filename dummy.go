package fhlbclient

import "github.com/valyala/fasthttp"

type DummyBalancer struct{}

func (b *DummyBalancer) Evaluate(_ []PenalizingClient) *PenalizingClient { return nil }

type DummyRequestHooks struct{}

func (h *DummyRequestHooks) PreRequest(_ *fasthttp.Request, _ *fasthttp.Response, _ *PenalizingClient) {
}
func (h *DummyRequestHooks) PostRequest(_ *fasthttp.Request, _ *fasthttp.Response, _ *PenalizingClient, _ error) {
}

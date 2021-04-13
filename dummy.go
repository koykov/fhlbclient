package fhlbclient

import "github.com/valyala/fasthttp"

type DummyRequestHooks struct{}

func (d *DummyRequestHooks) PreRequest(_ *fasthttp.Request, _ *fasthttp.Response, _ *PenalizingClient) {
}
func (d *DummyRequestHooks) PostRequest(_ *fasthttp.Request, _ *fasthttp.Response, _ *PenalizingClient, _ error) {
}

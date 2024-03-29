package fhlbclient

import "github.com/valyala/fasthttp"

// RequestHooker represents pre-/post-request hooks helper interface.
//
// Needs to perform some actions before and after request's execution.
type RequestHooker interface {
	PreRequest(req *fasthttp.Request, resp *fasthttp.Response, client *PenalizingClient)
	PostRequest(req *fasthttp.Request, resp *fasthttp.Response, client *PenalizingClient, err error)
}

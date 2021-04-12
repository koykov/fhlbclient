package fhlbclient

import "github.com/valyala/fasthttp"

// Health checker interface.
type HealthChecker interface {
	Check(req *fasthttp.Request, resp *fasthttp.Response, err error) bool
}

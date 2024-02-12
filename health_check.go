package fhlbclient

import "github.com/valyala/fasthttp"

// HealthChecker represents health checker interface.
type HealthChecker interface {
	Check(req *fasthttp.Request, resp *fasthttp.Response, err error) bool
}

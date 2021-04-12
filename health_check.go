package fhlbclient

import "github.com/valyala/fasthttp"

type HealthChecker interface {
	Check(req *fasthttp.Request, resp *fasthttp.Response, err error) bool
}

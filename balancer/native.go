package balancer

import (
	"github.com/koykov/fhlbclient"
)

// Native balancer implementation.
//
// Exactly reproduces native algorithm https://github.com/valyala/fasthttp/blob/master/lbclient.go#L99
// Main disadvantage: when any of clients becomes unavailable, LBClient thinks this client is the fastest and send all
// of requests to it. That's why this package exists.
type Native struct{}

func (b *Native) Evaluate(list []fhlbclient.PenalizingClient) *fhlbclient.PenalizingClient {
	minC := list[0]
	minN, minT := minC.RequestStats()
	for _, c := range list[1:] {
		n, t := c.RequestStats()
		if n < minN || (n == minN && t < minT) {
			minC = c
			minN = n
			minT = t
		}
	}
	return &minC
}

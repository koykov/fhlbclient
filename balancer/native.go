package balancer

import (
	"github.com/koykov/fhlbclient"
)

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

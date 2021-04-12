package balancer

import (
	"github.com/koykov/fhlbclient"
)

// Excluding balancer implementation.
//
// Is similar to Native balancer, but excludes clients under penalty from work till them become work as expected.
type Excluding struct{}

func (b *Excluding) Evaluate(list []fhlbclient.PenalizingClient) *fhlbclient.PenalizingClient {
	var (
		minC       *fhlbclient.PenalizingClient
		minN, minT uint64
	)
	for i := 0; i < len(list); i++ {
		if list[i].UnderPenalty() {
			continue
		} else {
			if minC == nil {
				minC = &list[i]
				minN, minT = minC.RequestStats()
				continue
			}
			pc := &list[i]
			n, t := pc.RequestStats()
			if n < minN || (n == minN && t < minT) {
				minC = pc
				minN = n
				minT = t
			}
		}
	}

	return minC
}

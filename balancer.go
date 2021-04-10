package fhlbclient

type Balancer interface {
	Evaluate([]PenalizingClient) *PenalizingClient
}

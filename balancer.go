package fhlbclient

// Balancer represents clients balancer interface.
type Balancer interface {
	Evaluate([]PenalizingClient) *PenalizingClient
}

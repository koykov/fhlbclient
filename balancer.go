package fhlbclient

// Clients balancer interface.
type Balancer interface {
	Evaluate([]PenalizingClient) *PenalizingClient
}

package fhlbclient

type MetricsWriter interface {
	// HostStatus registers count of status codes.
	HostStatus(host string, status int)
	// HostThrottle registers count of no alive clients errors.
	HostThrottle(host string)
}

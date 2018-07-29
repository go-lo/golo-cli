package main

// Localhost is a dummy LookerUpper which
// always returns localhost
type Localhost struct{}

// NewLocalhost will return an empty, dummy
// LookerUpper
func NewLocalhost() (l Localhost) {
	return
}

// Addresses will always return localhost
func (l Localhost) Addresses(_ string) (hb HostBinaryMap) {
	hb = make(HostBinaryMap)

	hb.AddHost("localhost")

	return
}

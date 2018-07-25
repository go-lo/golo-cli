package main

type Localhost struct{}

func NewLocalhost() (l Localhost) {
	return
}

func (l Localhost) Addresses(_ string) (hb HostBinaryMap) {
	hb = make(HostBinaryMap)

	hb.AddHost("localhost")

	return
}

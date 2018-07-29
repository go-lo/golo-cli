package main

import (
	"fmt"
	"os"
)

// HostBinary is a simple which puts a host against a binary
type HostBinary struct {
	Host   string
	Binary string
}

// HostBinaryMap is a map whose job is to just provide a simple
// way of searching for a HostBinary
type HostBinaryMap map[string]HostBinary

// AddHost adds a host to the HBM
func (hbm *HostBinaryMap) AddHost(host string) {
	(*hbm)[host] = HostBinary{Host: host}
}

// Add takes a HostBinary and puts it into the hbm
func (hbm *HostBinaryMap) Add(hb HostBinary) {
	m.Lock()
	(*hbm)[hb.Host] = hb
	m.Unlock()
}

// LookerUpper is an interface used by types which provide
// ways of finding agents by tag.
type LookerUpper interface {
	Addresses(tag string) HostBinaryMap
}

// SetLookerUpper takes a provider, from the cli, and tries to map
// it to a valid LookerUpper
func SetLookerUpper(provider string) (lu LookerUpper, err error) {
	switch provider {
	case "localhost":
		lu = NewLocalhost()

	case "digitalocean":
		t := os.Getenv("DO_TOKEN")

		lu, err = NewDigitalOcean(t)

	default:
		err = fmt.Errorf("No provider %q configured", *cloudProvider)
	}

	return
}

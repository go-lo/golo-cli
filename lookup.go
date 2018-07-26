package main

import (
	"fmt"
	"net/http"
	"os"
)

type HostBinary struct {
	Host   string
	Binary string
}

type HostBinaryMap map[string]HostBinary

func (hbm *HostBinaryMap) AddHost(host string) {
	(*hbm)[host] = HostBinary{}
}

func (hbm *HostBinaryMap) Add(hb HostBinary) {
	m.Lock()
	(*hbm)[hb.Host] = hb
	m.Unlock()
}

type LookerUpper interface {
	Addresses(string) HostBinaryMap
}

type oauthClient interface {
	Do(*http.Request) (*http.Response, error)
}

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

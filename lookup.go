package main

import (
	"net/http"
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

package main

import (
	"fmt"
	"strings"
)

// Env is a simple LookerUpper which
// will return addresses from the env
// var "GOLO_HOSTS=" - this is comma separated
type Env struct {
	hosts string
}

// NewEnv will return an empty LookerUpper
func NewEnv(hosts string) (e Env, err error) {
	if hosts == "" {
		err = fmt.Errorf("$GOLO_HOSTS is empty- have you set it?")

		return
	}

	e.hosts = hosts

	return
}

// Addresses will always return env
func (e Env) Addresses(_ string) (hb HostBinaryMap) {
	hb = make(HostBinaryMap)

	for _, h := range strings.Split(e.hosts, ",") {
		hb.AddHost(h)
	}

	return
}

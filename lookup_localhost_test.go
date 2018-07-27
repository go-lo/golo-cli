package main

import (
	"reflect"
	"testing"
)

func TestLocalhost_Addresses(t *testing.T) {
	l := NewLocalhost()

	expect := HostBinaryMap{"localhost": HostBinary{Host: "localhost"}}
	recvd := l.Addresses("")

	if !reflect.DeepEqual(expect, recvd) {
		t.Errorf("expected %+v, received %+v", expect, recvd)
	}
}

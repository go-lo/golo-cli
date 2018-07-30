package main

import (
	"testing"
)

func TestNewEnv(t *testing.T) {
	for _, test := range []struct {
		name        string
		hosts       string
		expectError bool
	}{
		{"happy path", "example.com", false},
		{"missing/ empty hosts", "", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			e, err := NewEnv(test.hosts)
			if test.expectError && err == nil {
				t.Errorf("expected error")
			}
			if !test.expectError && err != nil {
				t.Errorf("unexpected error %+v", err)
			}

			if test.hosts != e.hosts {
				t.Errorf("expected %q, received %q", test.hosts, e.hosts)
			}
		})
	}
}

func TestEnv_Addresses(t *testing.T) {
	e, _ := NewEnv("www1.example.com,www2.example.com")
	hbm := e.Addresses("")

	var ok bool

	_, ok = hbm["www1.example.com"]
	if !ok {
		t.Errorf("missing www1.example.com")
	}

	_, ok = hbm["www2.example.com"]
	if !ok {
		t.Errorf("missing www2.example.com")
	}
}

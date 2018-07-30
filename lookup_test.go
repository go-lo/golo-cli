package main

import (
	"os"
	"reflect"
	"testing"
)

func TestHostBinaryMap_AddHost(t *testing.T) {
	for _, test := range []struct {
		name   string
		hbm    HostBinaryMap
		hosts  []string
		expect HostBinaryMap
	}{
		{"Adding single host to empty HostBinaryMap", HostBinaryMap{}, []string{"1.example.com"}, HostBinaryMap{"1.example.com": HostBinary{Host: "1.example.com"}}},
		{"Adding multiple hosts to empty HostBinaryMap", HostBinaryMap{}, []string{"1.example.com", "2.example.com"}, HostBinaryMap{"1.example.com": HostBinary{Host: "1.example.com"}, "2.example.com": HostBinary{Host: "2.example.com"}}},
		{"Adding single host to non-empty HostBinaryMap", HostBinaryMap{"1.example.com": HostBinary{Host: "1.example.com"}}, []string{"2.example.com"}, HostBinaryMap{"1.example.com": HostBinary{Host: "1.example.com"}, "2.example.com": HostBinary{Host: "2.example.com"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, h := range test.hosts {
				test.hbm.AddHost(h)
			}

			if !reflect.DeepEqual(test.expect, test.hbm) {
				t.Errorf("expected %+v, received %+v", test.expect, test.hbm)
			}
		})
	}
}

func TestHostBinaryMap_Add(t *testing.T) {
	for _, test := range []struct {
		name   string
		hbm    HostBinaryMap
		hb     HostBinary
		expect HostBinaryMap
	}{
		{"Adding single host to empty HostBinaryMap", HostBinaryMap{}, HostBinary{"1.example.com", "abc123"}, HostBinaryMap{"1.example.com": HostBinary{Host: "1.example.com", Binary: "abc123"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.hbm.Add(test.hb)

			if !reflect.DeepEqual(test.expect, test.hbm) {
				t.Errorf("expected %+v, received %+v", test.expect, test.hbm)
			}
		})
	}
}

func TestSetLookerUpper(t *testing.T) {
	for _, test := range []struct {
		name        string
		provider    string
		env         map[string]string
		expectType  string
		expectError bool
	}{
		{"localhost", "localhost", make(map[string]string), "main.Localhost", false},
		{"digital ocean, correct env var", "digitalocean", map[string]string{"DO_TOKEN": "foo"}, "main.DigitalOcean", false},
		{"digital ocean, missing env var", "digitalocean", make(map[string]string), "main.DigitalOcean", true},
		{"env, correct var", "env", map[string]string{"GOLO_HOSTS": "example.com"}, "main.Env", false},
		{"env, missing env var", "env", make(map[string]string), "main.Env", true},

		{"No such provider", "non-such", make(map[string]string), "", true},
		{"emptyprovider", "", make(map[string]string), "", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			os.Clearenv()

			for k, v := range test.env {
				os.Setenv(k, v)
			}

			lu, err := SetLookerUpper(test.provider)

			if test.expectError && err == nil {
				t.Errorf("expected error")
			}

			if !test.expectError && err != nil {
				t.Errorf("unexpected error: %+v", err)
			}

			luType := reflect.TypeOf(lu)

			if luType != nil {
				if test.expectType != luType.String() {
					t.Errorf("expected type %q, received type %q", test.expectType, luType.String())
				}
			}

		})
	}
}

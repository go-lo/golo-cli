package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/digitalocean/godo"
)

type dummyDigitaloceanDropletService struct {
	ipv4 bool
	ipv6 bool
}

func (c dummyDigitaloceanDropletService) ListByTag(_ context.Context, _ string, _ *godo.ListOptions) (d []godo.Droplet, r *godo.Response, err error) {
	droplet := godo.Droplet{}

	if c.ipv4 || c.ipv6 {
		droplet.Networks = &godo.Networks{}
	}

	if c.ipv4 {
		droplet.Networks.V4 = []godo.NetworkV4{
			{IPAddress: "203.0.113.1"},
		}
	}

	if c.ipv6 {
		droplet.Networks.V6 = []godo.NetworkV6{
			{IPAddress: "2001:db8::"},
		}
	}

	d = append(d, droplet)

	return
}

func TestDigitalOcean_Addresses(t *testing.T) {
	for _, test := range []struct {
		name          string
		lookupService dummyDigitaloceanDropletService
		tag           string
		expect        HostBinaryMap
	}{
		{"Hosts with ipv6 addresses only", dummyDigitaloceanDropletService{ipv6: true}, "agentz", HostBinaryMap{"2001:db8::": HostBinary{"2001:db8::", ""}}},
		{"Hosts with ipv4 addresses only", dummyDigitaloceanDropletService{ipv4: true}, "agentz", HostBinaryMap{"203.0.113.1": HostBinary{"203.0.113.1", ""}}},
		{"Hosts with both types", dummyDigitaloceanDropletService{ipv4: true, ipv6: true}, "agentz", HostBinaryMap{"2001:db8::": HostBinary{"2001:db8::", ""}}},
		{"Hosts with no addresses", dummyDigitaloceanDropletService{}, "agentz", HostBinaryMap{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			d := DigitalOcean{
				lookupService: test.lookupService,
			}

			hbm := d.Addresses(test.tag)

			if !reflect.DeepEqual(test.expect, hbm) {
				t.Errorf("expected %+v, received %+v", test.expect, hbm)
			}
		})
	}
}

func TestDTS(t *testing.T) {
	dts := doTokenSource{
		AccessToken: "a-token",
	}

	tok, err := dts.Token()
	if err != nil {
		t.Errorf("unexpected error %+v", err)
	}

	if tok.AccessToken != dts.AccessToken {
		t.Errorf("oauth token value should equal %q, actually equals %q", dts.AccessToken, tok.AccessToken)
	}
}

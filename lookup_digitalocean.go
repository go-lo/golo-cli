package main

import (
	"context"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

type doTokenSource struct {
	AccessToken string
}

func (dts *doTokenSource) Token() (t *oauth2.Token, err error) {
	t = &oauth2.Token{
		AccessToken: dts.AccessToken,
	}

	return
}

type DigitalOcean struct {
	client *godo.Client
}

func NewDigitalOcean(token string) (d DigitalOcean, err error) {
	t := &doTokenSource{
		AccessToken: token,
	}

	d.client = godo.NewClient(oauth2.NewClient(context.Background(), t))

	return
}

func (do DigitalOcean) Addresses(tag string) (a HostBinaryMap) {
	droplets, _, _ := do.client.Droplets.ListByTag(context.TODO(), tag, nil)

	a = make(HostBinaryMap)

	for _, d := range droplets {
		n := d.Networks
		if n == nil {
			continue
		}

		if len(n.V6) > 0 {
			a.AddHost(n.V6[0].IPAddress)
		} else if len(n.V4) > 0 {
			a.AddHost(n.V4[0].IPAddress)
		}
	}

	return
}

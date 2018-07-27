package main

import (
	"context"
	"fmt"

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

type dropletLookupService interface {
	ListByTag(context.Context, string, *godo.ListOptions) ([]godo.Droplet, *godo.Response, error)
}

type DigitalOcean struct {
	lookupService dropletLookupService
}

func NewDigitalOcean(token string) (d DigitalOcean, err error) {
	if token == "" {
		err = fmt.Errorf("Missing digitalocean token- have you set $DO_TOKEN?")

		return
	}

	t := &doTokenSource{
		AccessToken: token,
	}

	d.lookupService = godo.NewClient(oauth2.NewClient(context.Background(), t)).Droplets

	return
}

func (do DigitalOcean) Addresses(tag string) (a HostBinaryMap) {
	droplets, _, _ := do.lookupService.ListByTag(context.TODO(), tag, nil)

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

package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type passingClient struct {
	body string
}

func (c passingClient) Do(_ *http.Request) (resp *http.Response, err error) {
	r := bytes.NewBufferString(c.body)
	resp = &http.Response{StatusCode: 201, Status: "ok", Body: ioutil.NopCloser(r)}

	return
}

type failingClient struct{}

func (c failingClient) Do(_ *http.Request) (resp *http.Response, err error) {
	r := bytes.NewBufferString("an error")
	resp = &http.Response{StatusCode: 500, Status: "error", Body: ioutil.NopCloser(r)}

	return
}

func TestQueueJob(t *testing.T) {
	for _, test := range []struct {
		name        string
		client      oauthClient
		hb          HostBinary
		j           Job
		expectError bool
	}{
		{"happy path", passingClient{`{"queued": true}`}, HostBinary{"example.com", "123abc"}, Job{}, false},
		{"error message from agent", passingClient{`Some agent error :(`}, HostBinary{"example.com", "123abc"}, Job{}, true},
		{"error message and status from agent", failingClient{}, HostBinary{"example.com", "123abc"}, Job{}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			client = test.client

			err := QueueJob(test.hb, test.j)
			if test.expectError && err == nil {
				t.Errorf("expected error")
			}

			if !test.expectError && err != nil {
				t.Errorf("unexpected error: %+v", err)
			}
		})
	}
}

func TestUploadSchedule(t *testing.T) {
	for _, test := range []struct {
		name        string
		client      oauthClient
		file        string
		address     string
		expect      HostBinary
		expectError bool
	}{
		{"happy path", passingClient{`{"binary": "abc123"}`}, "testdata/dummy-schedule", "example.com", HostBinary{"example.com", "abc123"}, false},
		{"error message from client", passingClient{`uh-oh :(`}, "testdata/dummy-schedule", "example.com", HostBinary{Host: "example.com"}, true},
		{"non-existent file", passingClient{`{"binary": "abc123"}`}, "testdata/nonsuch", "example.com", HostBinary{Host: "example.com"}, true},
		{"error from agent", failingClient{}, "testdata/dummy-schedule", "example.com", HostBinary{Host: "example.com"}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			client = test.client

			hb, err := UploadSchedule(test.file, test.address)
			if !reflect.DeepEqual(test.expect, hb) {
				t.Errorf("expected HostBinary %+v, received %+v", test.expect, hb)
			}

			if test.expectError && err == nil {
				t.Errorf("expected error")
			}

			if !test.expectError && err != nil {
				t.Errorf("unexpected error: %+v", err)
			}
		})
	}
}

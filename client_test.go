package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type passingClient struct {
	status int
	body   string
}

func (c passingClient) Do(_ *http.Request) (resp *http.Response, err error) {
	r := bytes.NewBufferString(c.body)
	resp = &http.Response{StatusCode: c.status, Status: "ok", Body: ioutil.NopCloser(r)}

	return
}

type failingClient struct{}

func (c failingClient) Do(_ *http.Request) (resp *http.Response, err error) {
	r := bytes.NewBufferString("an error")
	resp = &http.Response{StatusCode: 500, Status: "error", Body: ioutil.NopCloser(r)}

	return
}

type erroringClient struct{}

func (c erroringClient) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{}, fmt.Errorf("an error")
}

type conditionalClient struct {
	body    string
	errorOn string
}

func (c conditionalClient) Do(req *http.Request) (resp *http.Response, err error) {
	b := bytes.NewBufferString(c.body)
	s := 200

	if req.URL.Path == c.errorOn {
		b = bytes.NewBufferString("an error")
		s = 500
	}

	resp = &http.Response{StatusCode: s, Status: "", Body: ioutil.NopCloser(b)}

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
		{"happy path", passingClient{201, `{"queued": true}`}, HostBinary{"example.com", "123abc"}, Job{}, false},
		{"error message from agent", passingClient{201, `Some agent error :(`}, HostBinary{"example.com", "123abc"}, Job{}, true},
		{"error message and status from agent", failingClient{}, HostBinary{"example.com", "123abc"}, Job{}, true},
		{"erroring request", erroringClient{}, HostBinary{Host: "example.com"}, Job{}, true},
		{"Empty host", passingClient{201, `{"queued": true}`}, HostBinary{}, Job{}, true},
		{"response is 201, failed queue", passingClient{201, `{"queued": false}`}, HostBinary{Host: "example.com"}, Job{}, true},
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
		{"happy path", passingClient{200, `{"binary": "abc123"}`}, "testdata/dummy-schedule", "example.com", HostBinary{"example.com", "abc123"}, false},
		{"error message from client", passingClient{201, `uh-oh :(`}, "testdata/dummy-schedule", "example.com", HostBinary{Host: "example.com"}, true},
		{"non-existent file", passingClient{201, `{"binary": "abc123"}`}, "testdata/nonsuch", "example.com", HostBinary{Host: "example.com"}, true},
		{"error from agent", failingClient{}, "testdata/dummy-schedule", "example.com", HostBinary{Host: "example.com"}, true},

		{"erroring request", erroringClient{}, "testdata/dummy-schedule", "example.com", HostBinary{Host: "example.com"}, true},
		{"Empty host", passingClient{201, `{"binary": "abc123"}`}, "testdata/dummy-schedule", "", HostBinary{}, true},
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

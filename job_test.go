package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type passingEchoClient struct {
	hbm HostBinaryMap
}

func (c passingEchoClient) Do(r *http.Request) (resp *http.Response, err error) {
	hostParts := strings.Split(r.URL.Host, ":")

	hb, ok := c.hbm[hostParts[0]]
	if !ok {
		panic(fmt.Errorf("host %q not in %+v", hostParts[0], c.hbm))
	}

	b := bytes.NewBufferString(fmt.Sprintf(`{"binary": "%s"}`, hb.Binary))
	resp = &http.Response{StatusCode: 200, Status: "ok", Body: ioutil.NopCloser(b)}

	return
}

func TestUploadAndQueue(t *testing.T) {
	for _, test := range []struct {
		name        string
		job         Job
		client      oauthClient
		hbm         HostBinaryMap
		schedule    string
		expectError bool
	}{
		{"happy path", Job{}, passingClient{200, `{"queued": true, "binary": "abc123"}`}, HostBinaryMap{"example.com": HostBinary{"example.com", "abc123"}}, "testdata/dummy-schedule", false},
		{"error uploading schedule", Job{}, conditionalClient{`{"queued": true, "binary": "abc123"}`, "/upload"}, HostBinaryMap{"example.com": HostBinary{"example.com", "abc123"}}, "testdata/dummy-schedule", true},
		{"error queueing schedule", Job{}, conditionalClient{`{"queued": true, "binary": "abc123"}`, "/queue"}, HostBinaryMap{"example.com": HostBinary{"example.com", "abc123"}}, "testdata/dummy-schedule", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			client = test.client

			err := test.job.UploadAndQueue(test.hbm, test.schedule)

			if test.expectError && err == nil {
				t.Errorf("expected error")
			}

			if !test.expectError && err != nil {
				t.Errorf("unexpected error: %+v", err)
			}

		})
	}
}

func TestJob_Upload(t *testing.T) {
	// Ensure we can upload many schedules and not trip up against one another
	j := Job{
		Name:     "my-job",
		Users:    1024,
		Duration: 300,
	}

	hbm := HostBinaryMap{
		"www1.example.com": HostBinary{
			Host:   "www1.example.com",
			Binary: "binary-1",
		},
		"www2.example.com": HostBinary{
			Host:   "www2.example.com",
			Binary: "binary-2",
		},
		"www3.example.com": HostBinary{
			Host:   "www3.example.com",
			Binary: "binary-3",
		},
	}

	e := make(chan error, len(hbm))
	h := make(chan HostBinary, len(hbm))
	defer close(h)

	s := "testdata/job.toml"

	client = passingEchoClient{hbm}

	rcvdChan := make(chan int)
	go func() {
		for hb := range h {
			if hb.Binary == "" {
				t.Errorf("missing binary in %+v", hb)
			}

			if hbm[hb.Host].Binary != hb.Binary {
				t.Errorf("expected %q for host %q, received %q", hbm[hb.Host].Binary, hb.Host, hb.Binary)
			}

			rcvdChan <- 1
		}
	}()

	err := j.upload(hbm, s, h, e)
	if err != nil {
		t.Errorf("unexpected error %+v", err)
	}

	for i := 0; i < len(hbm); i++ {
		<-rcvdChan
	}
}

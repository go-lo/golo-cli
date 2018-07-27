package main

import (
	"testing"
)

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

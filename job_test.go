package main

import (
	"testing"
)

func TestUploadAndQueue(t *testing.T) {
	for _, test := range []struct {
		name        string
		job         Job
		hbm         HostBinaryMap
		schedule    string
		expectError bool
	}{
		{"happy path", Job{}, HostBinaryMap{"example.com": HostBinary{"example.com", "abc123"}}, "testdata/dummy-schedule", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			client = passingClient{`{"queued": true, "binary": "abc123"}`}

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

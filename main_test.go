package main

import (
	"testing"
)

type nop struct{}

func TestMain(t *testing.T) {
	for _, test := range []struct {
		name        string
		provider    string
		file        string
		schedule    string
		expectError bool
	}{
		{"happy path", "localhost", "testdata/job.toml", "testdata/dummy-schedule", false},
		{"missing file", "localhost", "testdata/non-such", "testdata/dummy-schedule", true},
		{"bad provider", "none", "testdata/job.toml", "testdata/dummy-schedule", true},
		{"misconfigured type map", "localhost", "testdata/invalid-config.toml", "testdata/dummy-schedule", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			oldVersionTypeMap := versionTypeMap
			defer func() {
				versionTypeMap = oldVersionTypeMap
			}()

			versionTypeMap["invalid"] = map[string]interface{}{
				latestVersion: new(nop),
			}

			client = passingClient{200, `{"queued": true, "binary": "abc123"}`}

			err := realmain(test.file, test.provider, "", test.schedule)
			if test.expectError && err == nil {
				t.Errorf("expected error")
			}

			if !test.expectError && err != nil {
				t.Errorf("unexpected error %+v", err)
			}
		})
	}
}

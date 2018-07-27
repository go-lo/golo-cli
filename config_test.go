package main

import (
	"reflect"
	"testing"
)

func TestTypeVersion(t *testing.T) {
	for _, test := range []struct {
		name          string
		inputType     string
		expectType    string
		expectVersion string
		expectError   bool
	}{
		{"Empty string", "", "", "", true},
		{"Malformed input", "config.v1.0.1", "", "", true},
		{"Type with missing version", "config", "config", latestVersion, false},
		{"Fully formed input type", "config.v1", "config", "v1", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			ty, v, err := typeVersion(test.inputType)

			if test.expectError && err == nil {
				t.Errorf("expected error")
			}

			if !test.expectError && err != nil {
				t.Errorf("unexpected error: %+v", err)
			}

			if test.expectType != ty {
				t.Errorf("expected type %q, received %q", test.expectType, ty)
			}

			if test.expectVersion != v {
				t.Errorf("expected type %q, received %q", test.expectVersion, v)
			}
		})
	}
}

func TestReadInput(t *testing.T) {
	for _, test := range []struct {
		name        string
		path        string
		expect      string
		expectError bool
	}{
		{"File exists, is valid", "testdata/job.toml", "*main.Job", false},
		{"Missing file", "testdata/nonsuch", "", true},
		{"File exists, is invalid config", "testdata/invalid-config.toml", "", true},
		{"File exists, is invalid toml", "testdata/invalid-toml.toml", "", true},
		{"File exists, valid type, wrong version", "testdata/job-wrong-version.toml", "", true},
		{"File exists, missing type/version", "testdata/job-no-type.toml", "", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			i, err := ReadInput(test.path)

			if test.expectError && err == nil {
				t.Errorf("expected error")
			}

			if !test.expectError && err != nil {
				t.Errorf("unexpected error: %+v", err)
			}

			iType := reflect.TypeOf(i)

			if iType != nil {
				if test.expect != iType.String() {
					t.Errorf("expected type %q, received type %q", test.expect, iType.String())
				}
			}
		})
	}
}

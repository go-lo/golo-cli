package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/mapstructure"
)

const (
	latestVersion = "latest"
)

var (
	versionTypeMap = map[string]map[string]interface{}{
		"job": {
			latestVersion: new(Job),
			"1":           new(Job),
		},
	}
)

type Input struct {
	Type   string      `yaml:"type"`
	Schema interface{} `yaml:"schema"`
}

func ReadInput(path string) (s interface{}, err error) {
	input, err := readConfig(path)
	if err != nil {
		return
	}

	i := new(Input)

	_, err = toml.Decode(string(input), &i)
	if err != nil {
		return
	}

	t, v, err := typeVersion(i.Type)
	if err != nil {
		return
	}

	tV, ok := versionTypeMap[t]
	if !ok {
		err = fmt.Errorf("Invalid type %q", t)

		return
	}

	s, ok = tV[v]
	if !ok {
		err = fmt.Errorf("No such version %q for type %q", v, t)

		return
	}

	err = mapstructure.Decode(i.Schema, s)

	return
}

func readConfig(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func typeVersion(s string) (t string, v string, err error) {
	if s == "" {
		err = fmt.Errorf("No type given")

		return
	}

	split := strings.Split(s, ".")

	switch len(split) {
	case 1:
		t = split[0]
		v = latestVersion

	case 2:
		t = split[0]
		v = split[1]

	default:
		err = fmt.Errorf("Invalid type")
	}

	return
}

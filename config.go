package main

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/mapstructure"
)

type Input struct {
	Type   string      `yaml:"type"`
	Schema interface{} `yaml:"schema"`
}

func ReadJobInput(path string) (j Job, err error) {
	input, err := readConfig(path)
	if err != nil {
		return
	}

	yI := new(Input)

	_, err = toml.Decode(string(input), &yI)
	if err != nil {
		return
	}

	if yI.Type != "job" {
		err = fmt.Errorf("Not a job")

		return
	}

	err = mapstructure.Decode(yI.Schema, &j)

	return
}

func readConfig(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

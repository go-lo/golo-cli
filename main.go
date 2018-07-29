package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	lu  LookerUpper
	hbm HostBinaryMap

	cloudProvider = flag.String("provider", "localhost", "Cloud Provider to query for agents")
	agentTag      = flag.String("agent-tag", "agent", "Tag value to query to find agents")
	schedule      = flag.String("schedule", "./schedule", "Schedule to upload")
	file          = flag.String("f", "config.toml", "Config file to load")

	m = sync.Mutex{}
)

func main() {
	flag.Parse()

	client = &http.Client{}

	err := realmain(*file, *cloudProvider, *agentTag, *schedule)
	if err != nil {
		log.Fatal(err)
	}
}

func realmain(f, p, t, s string) (err error) {
	i, err := ReadInput(f)
	if err != nil {
		return
	}

	lu, err = SetLookerUpper(p)
	if err != nil {
		return
	}

	log.Printf("Finding %s agents with tag %q", p, t)

	hbm = lu.Addresses(t)
	if len(hbm) == 0 {
		err = fmt.Errorf("No agents found")

		return
	}

	log.Printf("Found %d agents", len(hbm))

	switch i.(type) {
	case *Job:
		err = i.(*Job).UploadAndQueue(hbm, s)

	default:
		err = fmt.Errorf("No handler for %T", i)
	}

	return
}

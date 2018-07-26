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
	file          = flag.String("f", "config.yaml", "Config file to load")

	m = sync.Mutex{}
)

func main() {
	var err error
	flag.Parse()

	i, err := ReadInput(*file)
	if err != nil {
		log.Fatal(err)
	}

	lu, err = SetLookerUpper(*cloudProvider)
	if err != nil {
		log.Fatal(err)
	}

	client = &http.Client{}

	log.Printf("Finding %s agents with tag %q", *cloudProvider, *agentTag)

	hbm = lu.Addresses(*agentTag)

	log.Printf("Found %d agents", len(hbm))

	switch i.(type) {
	case *Job:
		err = UploadAndQueue(i.(*Job), hbm, *schedule)

	default:
		err = fmt.Errorf("No handler for %T", i)
	}

	if err != nil {
		log.Fatal(err)
	}
}

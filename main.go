package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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

	j, err := ReadJobInput(*file)
	if err != nil {
		log.Fatal(err)
	}

	switch *cloudProvider {
	case "localhost":
		lu = NewLocalhost()

	case "digitalocean":
		t := os.Getenv("DO_TOKEN")

		lu, err = NewDigitalOcean(t)
		if err != nil {
			panic(err)
		}

	default:
		log.Fatal(fmt.Errorf("No provider %q configured", *cloudProvider))
	}

	log.Printf("Finding %s agents with tag %q", *cloudProvider, *agentTag)

	hbm = lu.Addresses(*agentTag)

	log.Printf("Found %d agents", len(hbm))

	hostBinaries := make(chan HostBinary, len(hbm))
	errors := make(chan error, len(hbm))

	client = &http.Client{}

	for addr, _ := range hbm {
		go func(a string) {
			a = a

			log.Printf("%s - Starting Upload", a)
			hb, err := UploadSchedule(*schedule, a)

			log.Printf("%s - Completed Upload", a)

			hostBinaries <- hb
			errors <- err
		}(addr)
	}

	for i := 0; i < len(hbm); i++ {
		err := <-errors
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("All agents updated, queuing job")

	for i := 0; i < len(hbm); i++ {
		hb := <-hostBinaries
		go func(h HostBinary) {
			errors <- QueueJob(hb, j)
		}(hb)
	}

	for i := 0; i < len(hbm); i++ {
		err := <-errors
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Job successfully queued")
}

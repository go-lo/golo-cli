package main

import (
	"log"
)

type Job struct {
	Name     string `json:"name",yaml:"name"`
	Users    int    `json:"users",yaml:"users"`
	Duration int64  `json:"duration",yaml:"duration",`
	Binary   string `json:"binary"`
}

func UploadAndQueue(j *Job, hbm HostBinaryMap, schedule string) (err error) {
	hostBinaries := make(chan HostBinary, len(hbm))
	errors := make(chan error, len(hbm))

	log.Println("Updating agents")
	err = UploadJobs(hbm, schedule, hostBinaries, errors)
	if err != nil {
		return
	}

	log.Println("Queueing job")
	err = QueueJobs(j, hostBinaries, errors, len(hbm))
	if err != nil {
		return
	}

	return
}

func QueueJobs(j *Job, h chan HostBinary, errors chan error, size int) (err error) {
	for i := 0; i < size; i++ {
		hb := <-h

		go func(h HostBinary) {
			errors <- QueueJob(hb, *j)
		}(hb)
	}

	for i := 0; i < size; i++ {
		err = <-errors
		if err != nil {
			return
		}
	}

	return
}

func UploadJobs(hbm HostBinaryMap, schedule string, h chan HostBinary, errors chan error) (err error) {
	for addr, _ := range hbm {
		go func(a string) {
			a = a

			log.Printf("%s - Starting Upload", a)
			hb, err := UploadSchedule(schedule, a)

			log.Printf("%s - Completed Upload", a)

			h <- hb
			errors <- err
		}(addr)
	}

	for i := 0; i < len(hbm); i++ {
		err = <-errors
		if err != nil {
			return
		}
	}

	return
}

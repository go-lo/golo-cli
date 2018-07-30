[![Go Report Card](https://goreportcard.com/badge/github.com/go-lo/golo-cli)](https://goreportcard.com/report/github.com/go-lo/golo-cli)
[![Build Status](https://travis-ci.com/go-lo/golo-cli.svg?branch=master)](https://travis-ci.com/go-lo/golo-cli)
[![GoDoc](https://godoc.org/github.com/go-lo/golo-cli?status.svg)](https://godoc.org/github.com/go-lo/golo-cli)

# golo-cli

`golo-cli` is a CLI to interact with go-lo agents. It is able to:

 * Find agents by tags on DigitalOcean and when running on localhost
 * Upload schedule binaries
 * Queue jobs with varying durations, users, and anything else an agent accepts


## Installation

There is, currently, no packaged up `golo-cli` available. Instead it can be easily installed on machines with `go` installed. Should this not be the case, see [here](https://golang.org/doc/install)

```bash
$ go get github.com/go-lo/golo-cli
```

## Usage

All of the below examples assume a go-lo agent is running on localhost. See below for integration with cloud providers.

```bash
$ golo-cli --help
Usage of golo-cli:
  -agent-tag string
        Tag value to query to find agents (default "agent")
  -f string
        Config file to load (default "config.toml")
  -provider string
        Cloud Provider to query for agents (default "localhost")
  -schedule string
        Schedule to upload (default "./schedule")
```

### Running Loadtest Jobs

A job config file looks like:

```toml
# Formed of $type:$version. Omitting a version always defaults to 'latest'
type = "job:latest"

[schema]
name     = "my-loadtest"    # Name of the job
users    = 1024             # Users to simulate
duration = 300              # Length of time to run, in seconds
```

A job is paired with a `schedule`. A `schedule` is a binary which an agent runs; it _is_ the loadtest. See [go-lo/go-lo](https://github.com/go-lo/go-lo) for further information. **Note that this binary will need to be compiled with `CGO_ENABLED=0 GOOS=linux` set.**

Assuming our toml lives in `job.toml`, and our schedule `my-loadtest` we can start our job with:

```bash
$ golo-cli -f job.toml --schedule ./my-loadtest
```


## Integration with Cloud Providers

The CLI provides two flags pertinent to this topic:

```
  -agent-tag string
        Tag value to query to find agents (default "agent")
  -provider string
        Cloud Provider to query for agents (default "localhost")
```

Currently this software supports the following providers

| Provider     | Requirements                                        | Description                                                           |
|--------------|-----------------------------------------------------|-----------------------------------------------------------------------|
| localhost    | None                                                | Returns `localhost` and nothing else                                  |
| digitalocean | env var `$DO_TOKEN` containing a valid do API token | Will return 0 or more addresses for droplets containing specified tag |
| env          | `$GOLO_HOSTS` csv of hostnames of agents to use     | Returns the contents of `$GOLO_HOSTS` split on the comma              |


Providers are simple to write; a provider needs to implement

```golang
type LookerUpper interface {
    Addresses(tag string) HostBinaryMap
}
```

And then be registered in `lookup.go`

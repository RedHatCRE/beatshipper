# beatshipper

Since there’s no way to send GNU zip files through filebeat, this service will be responsible for checking if there are new .gz files based on a path that we’ll explode using globbing, decompress them, send them using the filebeat service with the provided configuration and store them in a local file registry.

It can also be used for non-compressed files.

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Requirements

> 🔔 Go Version 1.18 or higher

## Installation

- Move to the cloned repository and install with:

```
$ go install
```

```
$ beatshipper -h
Sends data based on paths that will be exploded using GLOB
with the possibility of passing GNU Zip files also.

Usage:
  beatshipper [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  send        Send the message as a beat

Flags:
  -h, --help      help for beatshipper
  -v, --version   version for beatshipper

Use "beatshipper [command] --help" for more information about a command.

$ beatshipper --version
beatshipper version 0.0.1
```

## Configuration

We should create the configuration in the following file `/etc/beatshipper/beatshipper-conf.yml`. An example is located in the root of the project. We should add the following fields:

- `host`
- `port`
- `path` array of paths that will be exploded using `GLOB`
- `registry`: the name of the file where we'll store the processed files.
- `recheck` time of recheck (re-process) if there are new pending files to send

## Usage:

If we have the right configuration, we can start the program using the following command - It will read the configuration, start checking if there are pending packages that aren't in the registry and send them to the server that is listening using `beats`- e.g. a logstash instance with the `beats` module -:

```
$ beatshipper send
2022/11/18 21:04:53 shipper.go:77: Conexion successful with: localhost:5044
2022/11/18 21:04:53 shipper.go:82: Processing: /jenkinst-test
jenkinst-test/jenkins-logs/rcj/test/sub-test/sub-sub-test/var/log/extra/rpm-list.txt.gz
2022/11/18 21:04:53 shipper.go:82: Processing: /jenkinst-test
jenkinst-test/jenkins-logs/rcj/test/sub-test/sub-sub-test10/var/log/extra/rpm-list.txt.gz
2022/11/18 21:04:53 shipper.go:118: Chunk slice into: 5 slices
2022/11/18 21:04:53 shipper.go:96: Sending batch of data...
2022/11/18 21:04:53 shipper.go:96: Sending batch of data...
2022/11/18 21:04:53 shipper.go:96: Sending batch of data...
2022/11/18 21:04:53 shipper.go:101: Conexion closed
2022/11/18 21:04:53 registry.go:75: Added to registry: /jenkinst-test
jenkinst-test/jenkins-logs/rcj/test/sub-test/sub-sub-test/var/log/extra/rpm-list.txt.gz
2022/11/18 21:04:53 registry.go:75: Added to registry: /jenkinst-test
jenkinst-test/jenkins-logs/rcj/test/sub-test/sub-sub-test10/var/log/extra/rpm-list.txt.gz
```

## Registry

Files will be stored in the following way:

```
$ jq < shipper_registry.json
{
  "ParsedFiles": [
    {
      "name": "/jenkinst-test/jenkins-logs/rcj/test/sub-test/sub-sub-test/var/log/extra/rpm-list.txt.gz",
      "date": "2022-11-18T21:04:53.773777134+01:00"
    },
    {
      "name": "/jenkinst-test/jenkins-logs/rcj/test/sub-test/sub-sub-test10/var/log/extra/rpm-list.txt.gz",
      "date": "2022-11-18T21:04:53.773781056+01:00"
    },
  ]
}
```
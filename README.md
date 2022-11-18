# beatshipper

Since there’s no way to send GNU zip files through filebeat, this service will be responsible for checking if there are new .gz files based on a path that we’ll explode using globbing, decompress them, send them using the filebeat service with the provided configuration and store them in a local file registry.

It can be used also for not compressed files.

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Requirements

[Download](https://golang.org/dl/) and install **Go**.

> 🔔 Please note: version 1.18 or higher required

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
```

## Configuration

We should create the configuration in the following file `/etc/beatshipper/beatshipper-conf.yml` located in the root of the project. We should add the following fields:

- `host`
- `port`
- `path` array of paths that will be exploded using `GLOB`
- `registry`: the name of the file where we'll store the processed files.
- `recheck` time of recheck (re-process)

## Usage:

- Manually:

```
$ beatshipper send
2022/11/18 21:04:53 shipper.go:77: Conexion successful with: localhost:5044
2022/11/18 21:04:53 shipper.go:82: Processing: /rhos-infra-dev-netapp/jenkins-logs/rcj/test/sub-test/sub-sub-test/var/log/extra/rpm-list.txt.gz
2022/11/18 21:04:53 shipper.go:82: Processing: /rhos-infra-dev-netapp/jenkins-logs/rcj/test/sub-test/sub-sub-test10/var/log/extra/rpm-list.txt.gz
2022/11/18 21:04:53 shipper.go:118: Chunk slice into: 5 slices
2022/11/18 21:04:53 shipper.go:96: Sending batch of data...
2022/11/18 21:04:53 shipper.go:96: Sending batch of data...
2022/11/18 21:04:53 shipper.go:96: Sending batch of data...
2022/11/18 21:04:53 shipper.go:101: Conexion closed
2022/11/18 21:04:53 registry.go:75: Added to registry: /rhos-infra-dev-netapp/jenkins-logs/rcj/test/sub-test/sub-sub-test/var/log/extra/rpm-list.txt.gz
2022/11/18 21:04:53 registry.go:75: Added to registry: /rhos-infra-dev-netapp/jenkins-logs/rcj/test/sub-test/sub-sub-test10/var/log/extra/rpm-list.txt.gz
```

## Registry

Files will be stored in the following way:

```
$ jq < shipper_registry.json
{
  "ParsedFiles": [
    {
      "name": "/rhos-infra-dev-netapp/jenkins-logs/rcj/test/sub-test/sub-sub-test/var/log/extra/rpm-list.txt.gz",
      "date": "2022-11-18T21:04:53.773777134+01:00"
    },
    {
      "name": "/rhos-infra-dev-netapp/jenkins-logs/rcj/test/sub-test/sub-sub-test10/var/log/extra/rpm-list.txt.gz",
      "date": "2022-11-18T21:04:53.773781056+01:00"
    },
  ]
}
```
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

We should create the configuration in any of the following both locations:

- `/etc/beatshipper/beatshipper-conf.yml`
- `$HOME/.config/beatshipper-conf.yml`

An example is located in the root of the project. We should add the following fields:

- `host`
- `port`
- `path` array of paths that will be exploded using `GLOB`
- `registry`: the name of the file where we'll store the processed files.
- `recheck` time of recheck (re-process) if there are new pending files to send
- `logsource` very useful in logstash to handle conditions based on the name of the source

# Configuration through SystemD

We have created two SystemD files:

- `service` file: it specifies the binary that will execute the service
- `timer` file: we are gonna execute the service as timer.

By default the timer will launch the service in one shoot type every 30 minutes (`*:0/30`). If we wanna change the time we should change the `OnCalendar` directime in the `timer` file.

These files are being coppied to the following directory `/lib/systemd/system/` if we install the generated RPM package.

We can enable and start the `timer` with the following commands:

```
$ systemctl enable beatshipper.timer
$ systemctl start beatshipper.timer
```

We can check if the timer has been activated:

```
$ systemctl list-timers | grep beatshipper -B1
NEXT                        LEFT          LAST                        PASSED       UNIT                           ACTIVATES
Mon 2022-12-26 16:00:00 CET 27min left    Mon 2022-12-26 15:30:32 CET 2min 12s ago beatshipper.timer              beatshipper.service

$ journalctl -u beatshipper.timer
dic 26 15:30:32 user systemd[1]: Started Beatshipper execution schedule.
```

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
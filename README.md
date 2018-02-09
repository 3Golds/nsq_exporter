# NSQ Exporter [![Build Status](https://travis-ci.org/timonwong/nsq_exporter.svg)][travis]

[![CircleCI](https://circleci.com/gh/timonwong/nsq_exporter/tree/master.svg?style=shield)][circleci]
[![Docker Repository on Quay](https://quay.io/repository/timonwong/nsq-exporter/status)][quay]
[![Docker Pulls](https://img.shields.io/docker/pulls/timonwong/nsq-exporter.svg?maxAge=604800)][hub]
[![Go Report Card](https://goreportcard.com/badge/github.com/timonwong/nsq_exporter)](https://goreportcard.com/report/github.com/timonwong/nsq_exporter)

NSQ exporter for prometheus.io, written in go.

## Building and running

### Build

```bash
make
```

### Running

```bash
./nsq_exporter <flags>
```

### Flags

Name                                       | Description
-------------------------------------------|--------------------------------------------------------------------------------------------------
--nsqd.addr                                | Address of the nsqd node. (default "http://localhost:4151/stats")
--timeout                                  | Timeout for trying to get stats from nsqd. (default 5s)
--collect                                  | Comma-separated list of collectors to use (available choices: `stats.topics`, `stats.channels` and `stats.clients`). (default "stats.topics,stats.channels")
--namespace                                | Namespace for the NSQ metrics. (default "nsq")
--log.level                                | Logging verbosity. (default: info)
--web.listen-address                       | Address to listen on for web interface and telemetry. (default: ":9118")
--web.telemetry-path                       | Path under which to expose metrics.
--version                                  | Print the version information.

## Using Docker

You can deploy this exporter using the Docker image from following registry:

* [DockerHub]\: [timonwong/nsq-exporter](https://registry.hub.docker.com/u/timonwong/nsq-exporter/)
* [Quay.io]\: [timonwong/nsq-exporter](https://quay.io/repository/timonwong/nsq-exporter)

For example:

```bash
docker pull timonwong/nsq-exporter

docker run -d -p 9117:9117 timonwong/nsq-exporter
```

[circleci]: https://circleci.com/gh/timonwong/nsq_exporter
[hub]: https://hub.docker.com/r/timonwong/nsq-exporter/
[travis]: https://travis-ci.org/timonwong/nsq_exporter
[quay]: https://quay.io/repository/timonwong/nsq-exporter
[DockerHub]: https://hub.docker.com
[Quay.io]: https://quay.io

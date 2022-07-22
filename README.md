# Tasmota Prometheus Exporter

[![Build Status](https://cloud.drone.io/api/badges/dyrkin/tasmota-exporter/status.svg?branch=master)](https://cloud.drone.io/dyrkin/luxmed-bot)
[![Docker Hub](https://img.shields.io/badge/image-latest-blue.svg?logo=docker&style=flat)](https://hub.docker.com/r/eugenezadyra/tasmota-exporter/tags/)

A Prometheus exporter for [Tasmota-enabled](https://tasmota.github.io/docs) devices.  
It subscribes to the topics `tele/+/+` and `stat/++` to which Tasmota devices send their updates.

![Grafana dashboard](grafana/dashboard.png)

Grafana dashboard is [available here](grafana/dashboard.json) on the GitHub repository.

## Prerequisites

* [Go](https://go.dev/doc/)

## Installation

### Build from the source code

You can pull the latest version of the source code and build it using the command line:

```bash
go mod download
go build -o tasmota-exporter ./cmd
chmod +x tasmota-exporter
```

Before running the application, you have to define the following environment variables:
```go
MQTT_HOSTNAME //default is localhost
MQTT_PORT //default is 1883
MQTT_USERNAME //default is empty
MQTT_PASSWORD //default is empty
PROMETHEUS_EXPORTER_PORT //listening port. Default is 9092
```

Then run it using:
```bash
./tasmota-exporter
```


### Using Docker

The exporter is also available as a [Docker image](https://hub.docker.com/r/eugenezadyra/tasmota-exporter/tags).
You can run it using the following example and pass configuration environment variables:

```bash
docker run \
  -e 'MQTT_HOSTNAME=192.168.1.10' \
  -e 'MQTT_PORT=1883' \
  -e 'MQTT_USERNAME=user' \
  -e 'MQTT_PASSWORD=password' \
  -e 'PROMETHEUS_EXPORTER_PORT=9092' \
  -p 9092:9092 \
  --restart unless-stopped \
  --name tasmota-exporter \
  eugenezadyra/tasmota-exporter:latest
```

### Using docker-compose

You can run it using docker-compose and the [template](compose/docker-compose.yml):

```yaml
version: '3.4'
services:
  tasmota-exporter:
    image: eugenezadyra/tasmota-exporter:latest
    environment:
      MQTT_HOSTNAME: 192.168.1.10
      MQTT_PORT: 1883
      MQTT_USERNAME: user
      MQTT_PASSWORD: password
      PROMETHEUS_EXPORTER_PORT: 9092
    ports:
      - '9092:9092'
    restart: unless-stopped
```

```bash
#run in the background
docker-compose up -d
#update to the latest version and restart
docker-compose pull
docker-compose restart
```

Once the exporter is running, you also have to update your `prometheus.yml` configuration to let it scrape the exporter:

```yaml
scrape_configs:
  - job_name:       'tasmota'
    scrape_interval: 5s
    static_configs:
      - targets: ['metrics.home:9092']
```

**!** ***Note:*** _Available metrics may vary by device_.
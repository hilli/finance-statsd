# finance-statsd

Collects finance statistics from Yahoo Finance and emits the data to StatsD (ie Telegraf/Datadog).

*THIS IS WORK IN PROGRESS!!!*

# Configuring

## Setup ticker symbols in the environment

```bash
docker run --rm -e SYMBOLS=AAPL,TSLA,MSFT hilli/finance-statsd
```

or with visual output in the log as well while collecting data at every 30s sending it to another server:

```bash
docker run --rm -e DEBUG=true -e COLLECTION_INTERVAL=30 -e STATSD_ENDPOINT=telegraf.local:8125 -e SYMBOLS=AAPL,TSLA,MSFT hilli/finance-statsd
```

If you prefer, you can create a `.env` file with the same environment as the docker-compose file and map it to the container as a volume.

## StatsD

You will need a `StatsD` server to accept your data. That could be [Telegraf](https://www.influxdata.com/time-series-platform/telegraf/), [Datadog](https://www.datadoghq.com), [Etsy's original StatsD](https://github.com/statsd/statsd) or anything that will accept StatsD formatted data. Set `STATSD_ENTRYPOINT` to point to your server of choice (hostname:port format).

# Docker Compose

`docker-compose.yaml` setting the environment:
```yaml
version: '3.7'
services:
  finance-statsd:
    image: hilli/finance-statsd
    container_name: finance-statsd
    hostname: finance-statsd
    restart: unless-stopped
    environment:
      STATSD_ENDPOINT: "localhost:8125"
      COLLECTION_INTERVAL: 60 # Seconds
      SYMBOLS: "AAPL,TSLA,MSFT"
      #DEBUG: "jearh" # Print the collected results in a human readable format to the docker log
```

`docker-compose.yaml` setting the environment in an `.env` file:
```yaml
version: '3.7'
services:
  finance-statsd:
    image: hilli/finance-statsd
    container_name: finance-statsd
    hostname: finance-statsd
    restart: unless-stopped
    volumes:
      - "./env:/.env"
```

`.env`:
```bash
STATSD_ENDPOINT=localhost:8125
COLLECTION_INTERVAL=60 # Seconds
SYMBOLS=AAPL,TSLA,MSFT
#DEBUG=OK # Output data to stdout as well
```

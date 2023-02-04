# URL Shortener Service

## Requirements
- It shortens the given URLs.
- It redirects to the original URL by getting a shortened URL.
- It provides metrics for monitoring.
- Scalability, Availability, Reliability

## System Design
TBD

## Sequence Diagram
### URL Shortening
TBD

### URL Redirection
TBD

## How to Run
### Option 1: Localhost
Install [redis](https://redis.io/docs/getting-started/installation/), [golang](https://go.dev/doc/install), and run:
```bash
$ redis-server
$ make run  # in another terminal
```

### Option 2: Docker
Install [docker](https://docs.docker.com/engine/install/) and run:
```bash
$ docker-compose up
```

### Option 3: Kubernetes
TBD

## API
```bash
POST /shorten  	# it returns a key value for shortened url
GET  /:key		# it redirects to the original url
GET  /docs		# swagger UI
GET  /metrics	# prometheus metrics
```

## Commands
```bash
make run            # build and run the project
make run-profile    # build and run the project with profiler
make setup-dev      # install go packages

# below commands are available after `make setup-dev`
make docs           # generate swagger ui
make format         # format the codes
make lint           # lint the codes
```

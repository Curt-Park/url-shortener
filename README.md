# URL Shortener Service

## Contents
- [Requirements](https://github.com/Curt-Park/url-shortener#requirements)
- [APIs](https://github.com/Curt-Park/url-shortener#apis)
- [System Design](https://github.com/Curt-Park/url-shortener#system-design)
  - Overview
  - URL Shortening
- [Sequence Diagram](https://github.com/Curt-Park/url-shortener#sequence-diagram)
  - URL Shortening
  - URL Redirection
- [How to Run](https://github.com/Curt-Park/url-shortener#how-to-run)
  - Host OS
  - Docker
  - Kubernetes
- [Test](https://github.com/Curt-Park/url-shortener#test)
  - Unit Test
  - Load Test
- [Tasks](https://github.com/Curt-Park/url-shortener#tasks)
- [Commands](https://github.com/Curt-Park/url-shortener#commands)

## Requirements
- It shortens the given URLs.
- It redirects to the original URL by getting a shortened URL.
- It provides metrics for monitoring.
- Scalability, Availability, Reliability.

## APIs
```bash
POST /shorten  # it returns a key value for shortened url
GET  /:key     # it redirects to the original url
GET  /docs     # swagger UI
GET  /metrics  # prometheus metrics
```

You can simply test it with `curl`.
```bash
$ curl -X 'POST' 'http://localhost:8080/shorten' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{ "url": "https://www.google.com/search?q=longlonglonglonglonglonglonglonglonglonglongurl" }'

{"key":"M8uIUx0W000"}
```

Go to http://localhost:8080/M8uIUx0W000 on your browser.
<img width="889" src="https://user-images.githubusercontent.com/14961526/216797605-61d64f76-0274-4dc5-a5c1-4df5aa23aca9.png">

## System Design
### Overview
![](https://user-images.githubusercontent.com/14961526/216781438-17cb9424-6239-4a37-94f0-14f18b0991c0.jpg)
- Server: [Echo](https://echo.labstack.com/) (Golang)
- Database: [Redis](https://redis.io/)

### URL Shortening
```mermaid
flowchart TD
  Start --> A
  A[Input: originalURL] --> B{Is it in DB?}
  B -->|Yes| C[Return the key for the short URL from DB]
  B -->|No| D[Generate an unique int64 value with snowflake]
  D --> E[Convert the unique key into a Base62 string]
  E --> F[Store the originalURL and the key]
  F --> C
  C --> End
```

## Sequence Diagram
### URL Shortening
```mermaid
sequenceDiagram
  autonumber
  actor U as User
  participant S as Server
  participant D as Database
  U ->> S: HTTP Req. POST Shortened URL {url}
  S ->> D: HTTP Req. GET Shortened URL {key}
  D -->> S: HTTP Resp. {key, exist}
  alt if not exists
    S ->> S: Generate Short URL key
    S ->> D: Store URL and key
  end
  S -->> U: HTTP Resp. Shortened URL key {key}
```

### URL Redirection
```mermaid
sequenceDiagram
  autonumber
  actor U as User
  participant S as Server
  participant D as Database
  U ->> S: HTTP Req. GET Original URL {key}
  S ->> D: HTTP Req. GET Original URL {key}
  D -->> S: HTTP Resp. {originalURL, exist}
  alt if not exists
    S -->> U: HTTP Resp. Not Found
  else
    S -->> U: HTTP Resp (Redirect). Found {originalURL}
  end
```

## How to Run
### Option 1: Host OS
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
Install [minikube]() and run:
```bash
make cluster  # init the k8s cluster
make charts   # install charts
```

To see grafana dashboard,
```bash
kubectl port-forward svc/prometheus-grafana 3000:80
```

Open http://localhost:3000/
- id: admin
- pw: prom-operator

## Test
### Unit Tests
```bash
make utest
```

### Load tests
You will need to install [Python3](https://www.python.org/downloads/) for this.
```bash
pip install locust  # just at the first beginning
make ltest
```

Open http://localhost:8089/

<img width="674" src="https://user-images.githubusercontent.com/14961526/216804990-87c9b65d-a150-482a-94f5-35e37ee00472.png">

## Tasks
- [x] APIs: url shortening, redirection, swagger UI, metrics
- [x] Code Formatting w/ `make format`
- [x] Code Linting w/ `make lint`
- [x] `Dockerfile` and `docker-compose.yaml`
- [x] Unit Test w/ [echo testing](https://echo.labstack.com/guide/testing/)
- [ ] Load Balancer (k8s)
- [ ] Auto Scaling (k8s)
- [ ] Ingress (k8s)
- [ ] SSL (k8s)
- [ ] Monitoring (k8s)
- [x] Load Tests w/ [Locust](https://locust.io/)

## Commands
```bash
make run            # build and run the project
make run-profile    # build and run the project with profiler
make setup-dev      # install go packages

# below commands are available after `make setup-dev`
make docs           # generate swagger ui
make format         # format the codes
make lint           # lint the codes

# tests
make utest          # run unit tests
make cover          # check the unit test coverage
make ltest          # load test w/ locust
```

# For users
run:
	$(MAKE) build
	./main

run-profile:
	$(MAKE) build
	./main --profile

# For devs
setup-dev:
	sh setup-dev.sh

format:
	golines -m 100 -t 4 -w main.go internal/*.go
	swag fmt

lint:
	golangci-lint run

# Swagger UI
.PHONY: docs
docs:
	swag init

# Tests
.PHONY: utest
utest:
	# Run `go help testflag` to see details
	go test -v -cover $(ARGS) ./internal

cover:
	ARGS="-coverprofile=cover.out" $(MAKE) utest
	go tool cover -html=cover.out

ltest:
	locust -f locustfile.py APIUser

# K8s Cluster
cluster:
	minikube start --driver=docker --extra-config=kubelet.housekeeping-interval=10s

.PHONY: charts
charts:
	# `helm uninstall name` for removal
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm dependency build charts/loki
	helm dependency build charts/promtail
	helm dependency build charts/prometheus
	helm dependency build charts/url-shortener
	helm dependency build charts/redis-cluster
	helm install loki charts/loki
	helm install promtail charts/promtail
	helm install prometheus charts/prometheus
	helm install redis charts/redis-cluster
	helm install url-shortener charts/url-shortener

finalize:
	minikube delete

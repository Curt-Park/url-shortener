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
	minikube addons enable metrics-server
	kubectl apply -f ingress
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update

.PHONY: charts
charts:
	# `helm uninstall name` for removal
	helm dependency build charts/loki
	helm dependency build charts/promtail
	helm dependency build charts/prometheus
	helm dependency build charts/url-shortener
	helm dependency build charts/redis
	helm install promtail charts/promtail
	helm install loki charts/loki
	helm install prometheus charts/prometheus
	helm install redis charts/redis
	helm install url-shortener charts/url-shortener

remove-charts:
	helm uninstall url-shortener || true
	helm uninstall redis || true
	helm uninstall prometheus || true
	helm uninstall loki || true
	helm uninstall promtail || true

finalize:
	minikube delete

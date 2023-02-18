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
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo add traefik https://traefik.github.io/charts
	helm repo add jetstack https://charts.jetstack.io
	helm repo update
	kubectl apply -f tls

.PHONY: charts
charts:
	# `helm uninstall name` for removal
	helm install cert-manager charts/cert-manager
	helm install traefik charts/traefik
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
	helm uninstall traefik || true
	helm uninstall cert-manager || true

finalize:
	minikube delete

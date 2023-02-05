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

.PHONY: docs
docs:
	swag init

format:
	golines -m 100 -t 4 -w main.go internal/*.go
	swag fmt

lint:
	golangci-lint run

.PHONY: utest
utest:
	# Run `go help testflag` to see details
	go test -v -cover $(ARGS) ./internal

cover:
	ARGS="-coverprofile=cover.out" $(MAKE) utest
	go tool cover -html=cover.out

ltest:
	locust -f locustfile.py APIUser

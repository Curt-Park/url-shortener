run:
	$(MAKE) build
	./main

run-profile:
	$(MAKE) build
	./main --profile

docs:
	swag init

setup-dev:
	sh setup-dev.sh

format:
	golines -m 100 -t 4 -w main.go internal/user/*.go
	# swag fmt

lint:
	golangci-lint run

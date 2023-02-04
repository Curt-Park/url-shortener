#!/bin/sh
go install github.com/swaggo/swag/cmd/swag@v1.8.7  # swagger cli
go install golang.org/x/tools/cmd/goimports@v0.3.0  # format
go install github.com/segmentio/golines@v0.11.0  # format
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1  # lint

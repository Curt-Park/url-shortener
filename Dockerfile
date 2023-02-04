FROM golang:1.19.3-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /build
COPY . .
RUN go build main.go

FROM scratch
COPY --from=builder /build/main .
ENTRYPOINT ["/main"]

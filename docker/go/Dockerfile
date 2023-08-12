FROM golang:1.21.0-bookworm AS go-builder

WORKDIR /go-ddd-example

COPY . .
COPY start.sh /scripts/start.sh

RUN go install github.com/rakyll/gotest@latest
RUN go install github.com/cosmtrek/air@latest

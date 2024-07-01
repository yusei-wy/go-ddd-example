FROM golang:1.22

WORKDIR /app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/houqp/sqlvet@latest
RUN go install github.com/kyoh86/richgo@latest
RUN go install github.com/lib/pq@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

COPY . .

# go_ddd_example-server

## Requirements

- Docker
- [golang-migrate](https://github.com/golang-migrate/migrate)

## Directories

This project is a sample web application created using Go, employing the DDD, CQRS, and Feature Sliced Design (FSD) patterns.

```
.
├── app
│   └── env
├── cmd
├── db
│   └── migrations
├── docker
│   └── app
├── feature
│   └── <domain name>
│       ├── domain
│       │   └── model
│       │       └── command
│       │       └── query
│       │       └── value_object
│       ├── infrastructure
│       │       └── in_memory_repository
│       │       └── psql_repository
│       └── usecase
├── server
├── share
│   ├── custom_error
│   ├── domain
│   │   └── model
│   │       └── value_object
│   └── usecase
└── tmp
```

## commands

```
| Command  | Description |
| --- | --- |
| `make dev` | run dev server |
| `make reset-db` | clean db & create db & migrate |
| `make lint` | run golangci-lint |
| `make format` | run golangci-lint --fix |
| `make test` | run go test ./... |
```


## TODO

- [x] DDD architecture
- [x] Custom error handler
- [x] Create getters directive
- [ ] Create new directive
- [x] Migration system
- [ ] Generate Request And Response Type With OpenAPI
- [x] Dependency Injection

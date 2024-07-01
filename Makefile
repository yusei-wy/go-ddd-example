DB_CONTAINER=postgres
DB_NAME=go_ddd_example

include .env.dev

init:
	# set up git hooks
	ln -sf $(CURDIR)/pre-commit $(CURDIR)/.git/hooks/pre-commit

clean-db:
	docker-compose exec -it $(DB_CONTAINER) psql -U postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"

create-db:
	docker-compose exec -it $(DB_CONTAINER) psql -U postgres -c "CREATE DATABASE $(DB_NAME);"

migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(NAME)

reset-db: clean-db create-db migrate-up

migrate-up:
	migrate -database $(DATABASE_URL) -path db/migrations up

migrate-down:
	migrate -database $(DATABASE_URL) -path db/migrations down

dump: reset-db
	docker-compose exec -it $(DB_CONTAINER) pg_dump --schema-only -U postgres $(DB_NAME) > db/dump.sql

dev: reset-db
	lsof -t -i tcp:8080 | xargs kill -9
	DATABASE_URL=$(DATABASE_URL) air

lint:
	golangci-lint run ./...

format:
	gofumpt -e -d -l -w .
	golangci-lint run --fix ./...

test:
	go test -v ./...

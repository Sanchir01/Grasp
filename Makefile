PHONY:
SILENT:
MIGRATION_NAME ?= new_migration

build:
	go build -o ./.bin/main ./cmd/main/main.go

run: build
	./.bin/main

build-image:
	docker build -t cryptobot-dockerfile .
start-container:
	docker run --name cryptobot-test -p 80:80 --env-file .env cryptobot-dockerfile
migrations-up:
	goose -dir internal/db/migrations postgres "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable" up

migrations-status:
	goose -dir internal/db/migrations postgres "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable" status

migrations-new:
	goose -dir internal/db/migrations create $(MIGRATION_NAME) sql

compose-up:
	docker-compose -f docker-compose.yaml ps

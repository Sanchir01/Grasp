PHONY:
SILENT:

build:
	go build -o ./.bin/main ./cmd/main/main.go

run: build
	./.bin/main

build-image:
	docker build -t cryptobot-dockerfile .
start-container:
	docker run --name cryptobot-test -p 80:80 --env-file .env cryptobot-dockerfile


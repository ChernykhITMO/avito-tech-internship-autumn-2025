APP_NAME=avito-pr-service
CMD_PATH=./cmd/app

.PHONY: build run test lint docker-up docker-down migrate

build:
	go build -o bin/$(APP_NAME) $(CMD_PATH)

run:
	go run $(CMD_PATH)

test:
	go test ./...

lint:
	golangci-lint run ./...

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down


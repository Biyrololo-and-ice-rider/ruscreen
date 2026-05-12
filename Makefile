.PHONY: build run dev test test-unit test-integration lint docker-up docker-down generate

build:
	go build -o bin/meet ./cmd/server

run: build
	./bin/meet

dev:
	docker-compose up -d && go run ./cmd/server

test-unit:
	go test ./internal/domain/... ./internal/usecase/... -v -race -count=1

test-integration:
	go test ./internal/infra/... ./internal/adapter/... -v -race -count=1 -tags=integration

test: test-unit test-integration

lint:
	golangci-lint run ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

generate:
	go generate ./...

BIN := "./bin/calendar"
SCHEDULER_BIN := "./bin/calendar_scheduler"
SENDER_BIN := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build-calendar:
	GOOS=linux go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-scheduler:
	GOOS=linux go build -v -o $(SCHEDULER_BIN) ./cmd/scheduler

build-sender:
	GOOS=linux go build -v -o $(SENDER_BIN) ./cmd/sender

build: build-calendar build-scheduler build-sender

build-win-calendar:
	GOOS=windows GOARCH=amd64 go build -v -o $(BIN).exe -ldflags "$(LDFLAGS)" ./cmd/calendar

build-win-scheduler:
	GOOS=windows GOARCH=amd64 go build -v -o $(SCHEDULER_BIN).exe ./cmd/scheduler

build-win-sender:
	GOOS=windows GOARCH=amd64 go build -v -o $(SENDER_BIN).exe ./cmd/sender

build-win: build-win-calendar build-win-scheduler build-win-sender

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f deployments/images/api/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./cmd/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

lint: install-lint-deps
	golangci-lint run ./...

generate:
	protoc ./internal/server/grpc/calendar.proto --go_out=./internal/server/grpc/pb --go-grpc_out=./internal/server/grpc/pb

up: 
	docker compose up -d --build

down:
	docker compose down

run-integration-test:
	go test ./intergation_tests/...

integration-tests: up run-integration-test down

.PHONY: build run build-img run-img version test lint

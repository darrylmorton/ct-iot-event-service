.DEFAULT_GOAL := build
BINARY_NAME=server.out
EVENTS_DB_DSN_VALUE="postgres://postgres:postgres@localhost:5432/events?sslmode=disable"

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./cmd/api
.PHONY:vet

build: vet
	go build -v -o ${BINARY_NAME} ./cmd/api
.PHONY:build

run-dev:
	go build -v -o ${BINARY_NAME} ./cmd/api
	EVENTS_DB_DSN=${EVENTS_DB_DSN_VALUE} ./${BINARY_NAME}
.PHONY:run-dev

run-prod:
	go build -v -o ${BINARY_NAME} ./cmd/api
	./${BINARY_NAME}
.PHONY:run-prod

clean:
	go clean
	rm ${BINARY_NAME}
.PHONY:clean

test:
	go test -v ./...
.PHONY:test

build:
	go build -o calendar cmd/main.go

test:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./...

install:
	go mod download

install-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

lint: install-deps
	golangci-lint run ./...

migrate:
	goose -dir migrations/ postgres "host=localhost port=5432 user=otus password=otus dbname=calendar sslmode=disable" up

.PHONY: build

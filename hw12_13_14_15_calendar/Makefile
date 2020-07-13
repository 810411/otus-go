build:
	go build -o calendar cmd/calendar/main.go
	go build -o calendar_sender ./cmd/calendar_sender/main.go
	go build -o calendar_scheduler ./cmd/calendar_scheduler/main.go

test:
	go test -v -count=1 -race -timeout=30s ./...

integration-tests:
	set -e ;\
    docker-compose -f ./deployments/docker-compose.test.yml up --build -d ;\
    test_status_code=0 ;\
    docker-compose -f ./deployments/docker-compose.test.yml run integration_tests go test -v || test_status_code=$$? ;\
    docker-compose -f ./deployments/docker-compose.test.yml down ;\
    exit $$test_status_code ;\

install:
	go mod download

install-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

lint: install-deps
	golangci-lint run ./...

migrate:
	goose -dir migrations/ postgres "host=localhost port=5432 user=otus password=otus dbname=calendar sslmode=disable" up

generate:
	go generate ./...

up:
	docker-compose -f ./deployments/docker-compose.yml up

down:
	docker-compose -f ./deployments/docker-compose.yml down

.PHONY: build

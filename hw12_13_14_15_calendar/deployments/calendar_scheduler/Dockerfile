FROM golang:1.14-alpine AS builder
LABEL stage=builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o calendar_scheduler cmd/calendar_scheduler/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/calendar_scheduler /app/calendar_scheduler

ENTRYPOINT ["./calendar_scheduler"]

FROM golang:1.14-alpine AS builder
LABEL stage=builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o calendar_sender cmd/calendar_sender/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/calendar_sender /app/calendar_sender

ENTRYPOINT ["./calendar_sender"]

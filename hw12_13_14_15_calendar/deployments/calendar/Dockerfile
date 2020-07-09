FROM golang:1.14-alpine AS builder
LABEL stage=builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o calendar cmd/calendar/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/calendar /app/calendar

EXPOSE 8888 50051

ENTRYPOINT ["./calendar"]

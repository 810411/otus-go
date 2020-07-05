FROM golang:1.14-alpine AS builder
LABEL stage=builder

RUN set -ex \
        && apk add --no-cache --virtual build-dependencies \
            build-base \
            git \
        && go get -ldflags "-extldflags -static" github.com/pressly/goose/cmd/goose \
        && apk del build-dependencies

FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/bin/goose /app/goose

ENTRYPOINT ["./goose"]

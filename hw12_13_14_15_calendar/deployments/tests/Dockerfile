FROM golang:1.14

WORKDIR /tests

COPY ./tests/integration .

RUN go mod download

CMD ["go", "test", "-v"]

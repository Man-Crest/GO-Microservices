FROM golang:latest AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

CMD ["/app/brokerApp"]


# build i smaller docker image
FROM alpine:latest

RUN mkdir /app

COPY brokerApp /app

CMD ["/app/brokerApp"]





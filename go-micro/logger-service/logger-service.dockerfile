FROM golang:latest AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o loggerServiceApp ./cmd/api

CMD ["/app/loggerServiceApp"]

FROM alpine:latest

RUN mkdir /app

COPY loggerServiceApp /app

CMD [ "/app/loggerServiceApp"]
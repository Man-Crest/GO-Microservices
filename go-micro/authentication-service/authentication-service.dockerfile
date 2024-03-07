# FROM golang:latest AS builder

# RUN mkdir /app

# COPY . /app

# WORKDIR /app

# RUN CGO_ENABLED=0 go build -o authApp ./cmd/api

# CMD ["/app/authApp"]

# build i smaller docker image
FROM alpine:latest

RUN mkdir /app

COPY authApp /app

CMD [ "/app/authApp"]





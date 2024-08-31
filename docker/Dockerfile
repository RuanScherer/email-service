FROM golang:1.22 AS build

WORKDIR /app

COPY . .

RUN go build -o email-service

FROM ubuntu:latest

COPY --from=build /app/email-service .

RUN apt-get update \
  && apt-get install build-essential librdkafka-dev -y \
  && apt-get clean

ENTRYPOINT ["./email-service"]

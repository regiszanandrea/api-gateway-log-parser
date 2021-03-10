FROM golang:alpine

LABEL version="1/0"
LABEL description="This image provides a golang environment with AWS cli"
MAINTAINER regiszanandrea@gmail.com

COPY . /app
WORKDIR /app

RUN go env -w GO111MODULE=auto
RUN go mod tidy
RUN env GOOS=linux go build -o bin `go list ./cmd/...`
RUN chmod -R +x cmd

RUN env GOOS=linux go build -o bin/migrations `go list ./db/...`
RUN chmod -R +x db

ENTRYPOINT ["tail", "-f", "/dev/null"]
## AWS API Gateway Log Parser

### Dependencies

- Docker
- Docker Compose

### Optional Dependencies

- [Golang](https://golang.org/doc/install)

### Installation

1. Start docker containers

```sh
docker-compose up -d
```

2. Add the Git hooks to your local .git directory

```sh
cp -a githooks/* .git/hooks
```

## Running the application

### Local Development

Todo

### Commands

See Makefile for available commands!

## Code Architecture

[See diagram here](https://confluence.endurance.com/display/LED/1.+Golang+Microservice)

### File structure

```
├── application
│   ├── handler
│   └── healthcheck # Healthcheck business domain
│     ├── model.go
│     └── services.go
├── cmd  # All inputs to the application
│   └── healthcheck
│       └── main.go
├── internal # Business domains
│   ├── di # Dependecy injection
│   │   ├── container.go
├── githooks
│   ├── pre-commit
│   └── pre-push
├── docker-compose.yml
├── Makefile
├── go.mod

```
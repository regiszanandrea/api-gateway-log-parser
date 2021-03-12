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

2. Create table

```sh
make migrate
```

3. Put the log file on `assets` folder, this folder will be visible on docker container, then execute the follow command
   to parse the file:

```sh
make FILE_PATH=/data/{fileName} parse
```

2. Add the Git hooks to your local .git directory

```sh
cp -a githooks/* .git/hooks
```

3. Install Git hooks

```sh
make install-hooks
```

## Running the application

### Commands

See Makefile for available commands!

e.g.:
To generate CSV file by service

```
make SERVICE=c3e86413-648a-3552-90c3-b13491ee07d6 export-by-service
```

All files generated will be on `assets` folder

## Testing

To test the application, there are some commands on Makefile:

- Generate coverage: this will be open your browser with the code coverage detailed

```
make generate-coverage
```


## Code Architecture

All code folders is based on [project-layout github](https://github.com/golang-standards/project-layout) and
use [Hexagonal architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)) pattern:

![Hexagonal Architecture](https://camo.githubusercontent.com/5b802633e416330b3b1c3e55df695291680b60b074d8341cc76199ec30d1707e/68747470733a2f2f692e696d6775722e636f6d2f65736557566c422e706e67)

### File structure

```
├── application
│   ├── handler
│   │   ├── export_by_consumer.go
│   │   ├── export_by_service.go
│   │   ├── export_metrics_by_service.go
│   │   └── log_parser_handler.go
│   └── service
│       ├── agigateway_integration_test.go
│       ├── apigateway.go
│       └── apigateway_test.go
├── assets
├── bin
│   └── migrations
├── cmd
│   ├── apigateway_log_parser
│   │   └── main.go
│   ├── export_by_consumer
│   │   └── main.go
│   ├── export_by_service
│   │   └── main.go
│   └── export_metrics_by_service
│       └── main.go
├── data
├── db
│   └── dynamodb
│       └── migrations
│           └── create_logs_table
│               └── create_table.go
├── docker-compose.yml
├── Dockerfile
├── githooks
│   ├── pre-commit
│   └── pre-push
├── go.mod
├── go.sum
├── internal
│   └── di
│       └── container.go
├── LICENSE
├── Makefile
├── pkg
│   ├── apigateway
│   │   ├── apigateway.go
│   │   └── repository
│   │       ├── driver
│   │       │   ├── driver.go
│   │       │   └── dynamodb.go
│   │       └── repository.go
│   └── filesystem
│       ├── filesystem.go
│       └── local.go
├── README.md
├── test
│   ├── handler
│   │   ├── export_by_consumer_integration_test.go
│   │   ├── export_by_service_integration_test.go
│   │   ├── export_metrics_by_service_integration_test.go
│   │   └── log_parser_handler_integration_test.go
│   └── mocks
│       ├── driver.go
│       └── filesystem.go
```
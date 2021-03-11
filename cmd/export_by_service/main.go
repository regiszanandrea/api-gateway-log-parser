package main

import (
	"api-gateway-log-parser/internal/di"
	"context"
	"log"
)

var container *di.Container

func init() {
	container = di.NewContainer()
}

func main() {

	handle := container.GetExportByServiceHandler()

	err := handle(context.Background())

	if err != nil {
		log.Fatal(err)
	}
}

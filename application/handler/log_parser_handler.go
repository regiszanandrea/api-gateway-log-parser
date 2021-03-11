package handler

import (
	"api-gateway-log-parser/pkg/apigateway"
	"context"
	"errors"
	"os"
)

type LogParserHandler struct {
	service apigateway.LogService
}

func NewLogParserHandler(service apigateway.LogService) *LogParserHandler {
	return &LogParserHandler{service: service}
}

func (h *LogParserHandler) HandleApiGatewayLogParser(ctx context.Context) error {
	if len(os.Args) < 2 {
		return errors.New("path parameter not provided")
	}

	path := os.Args[1]

	if path == "" {
		return errors.New("path parameter could not be empty")
	}

	return h.service.Parse(path)
}

package handler

import (
	"api-gateway-log-parser/pkg/apigateway"
	"context"
	"errors"
	"os"
)

type ExportByConsumerHandler struct {
	service apigateway.LogService
}

func NewExportByConsumerHandler(service apigateway.LogService) *ExportByConsumerHandler {
	return &ExportByConsumerHandler{service: service}
}

func (h *ExportByConsumerHandler) HandleExportByConsumer(ctx context.Context) error {
	if len(os.Args) < 2 {
		return errors.New("consumer parameter not provided")
	}

	consumer := os.Args[1]

	if consumer == "" {
		return errors.New("consumer parameter could not be empty")
	}

	return h.service.ExportByConsumer(consumer)
}

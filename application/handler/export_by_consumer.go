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

var (
	ErrConsumerParameterNotFound        = errors.New("service parameter not provided")
	ErrConsumerParameterCouldNotBeEmpty = errors.New("service parameter could not be empty")
)

func NewExportByConsumerHandler(service apigateway.LogService) *ExportByConsumerHandler {
	return &ExportByConsumerHandler{service: service}
}

func (h *ExportByConsumerHandler) HandleExportByConsumer(ctx context.Context) error {
	if len(os.Args) < 2 {
		return ErrConsumerParameterNotFound
	}

	consumer := os.Args[1]

	if consumer == "" {
		return ErrConsumerParameterCouldNotBeEmpty
	}

	return h.service.ExportByConsumer(consumer)
}

package handler

import (
	"api-gateway-log-parser/pkg/apigateway"
	"context"
	"errors"
	"os"
)

type ExportByServiceHandler struct {
	service apigateway.LogService
}

func NewExportByServiceHandler(service apigateway.LogService) *ExportByServiceHandler {
	return &ExportByServiceHandler{service: service}
}

var (
	ErrServiceParameterNotFound        = errors.New("service parameter not provided")
	ErrServiceParameterCouldNotBeEmpty = errors.New("service parameter could not be empty")
)

func (h *ExportByServiceHandler) HandleExportByService(ctx context.Context) error {
	if len(os.Args) < 2 {
		return ErrServiceParameterNotFound
	}

	service := os.Args[1]

	if service == "" {
		return ErrServiceParameterCouldNotBeEmpty
	}

	return h.service.ExportByService(service)
}

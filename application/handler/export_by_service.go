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

func (h *ExportByServiceHandler) HandleExportByService(ctx context.Context) error {
	if len(os.Args) < 2 {
		return errors.New("service parameter not provided")
	}

	service := os.Args[1]

	if service == "" {
		return errors.New("service parameter could not be empty")
	}

	return h.service.ExportByService(service)
}

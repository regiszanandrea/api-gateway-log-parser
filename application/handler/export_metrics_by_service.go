package handler

import (
	"api-gateway-log-parser/pkg/apigateway"
	"context"
	"errors"
	"os"
)

type ExportMetricsByServiceHandler struct {
	service apigateway.LogService
}

func NewExportMetricsByServiceHandler(service apigateway.LogService) *ExportMetricsByServiceHandler {
	return &ExportMetricsByServiceHandler{service: service}
}

func (h *ExportMetricsByServiceHandler) HandleExportMetricsByService(ctx context.Context) error {
	if len(os.Args) < 2 {
		return errors.New("service parameter not provided")
	}

	service := os.Args[1]

	if service == "" {
		return errors.New("service parameter could not be empty")
	}

	return h.service.ExportMetricsByService(service)
}

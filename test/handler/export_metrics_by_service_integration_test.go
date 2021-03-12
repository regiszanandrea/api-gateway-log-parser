// +build integration

package test

import (
	"api-gateway-log-parser/application/handler"
	"api-gateway-log-parser/application/service"
	"api-gateway-log-parser/pkg/apigateway"
	"api-gateway-log-parser/pkg/apigateway/repository"
	mock "api-gateway-log-parser/test/mocks"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	as "github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"
	"os"
	"testing"
)

func TestHandleExportMetricsByService_ShouldReturnErrorWithWrongParameters(t *testing.T) {
	assert := as.New(t)

	filesystem := mock.FileSystemMock{}
	driverMock := mock.DriverMock{}

	repo := repository.NewApiGatewayLogRepository(&driverMock)

	s, _ := service.NewApiGatewayLogParserService(repo, &filesystem)

	h := handler.NewExportMetricsByServiceHandler(s)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{}

	err := h.HandleExportMetricsByService(context.Background())

	assert.NotNil(err)
	assert.Same(err, handler.ErrServiceParameterNotFound)

	os.Args = []string{"", ""}

	err = h.HandleExportMetricsByService(context.Background())

	assert.NotNil(err)
	assert.Same(err, handler.ErrServiceParameterCouldNotBeEmpty)
}

func TestHandleExportMetricsByService_ShouldExportMetrics(t *testing.T) {
	assert := as.New(t)

	var logs []*apigateway.Log

	consumerID := "29a5a16b-e4fa-331f-9f1c-5adea563d7de"
	serviceID := "c3e86413-648a-3552-90c3-b13491ee07d6"
	userIP := "0.0.0.0"

	itemsPerPage := 1000

	logs = append(logs, &apigateway.Log{
		Request:             apigateway.Request{},
		UpstreamURI:         "/",
		Response:            apigateway.Response{},
		AuthenticatedEntity: apigateway.AuthenticatedEntity{ConsumerID: apigateway.Consumer{UUID: consumerID}},
		Route:               apigateway.Route{},
		Service:             apigateway.Service{},
		Latencies: apigateway.Latencies{
			Proxy:   1,
			Gateway: 2,
			Request: 3,
		},
		ClientIP:   userIP,
		StartedAt:  12345,
		ServiceID:  serviceID,
		ConsumerID: consumerID,
	})

	logs = append(logs, &apigateway.Log{
		Request:             apigateway.Request{},
		UpstreamURI:         "/",
		Response:            apigateway.Response{},
		AuthenticatedEntity: apigateway.AuthenticatedEntity{ConsumerID: apigateway.Consumer{UUID: consumerID}},
		Route:               apigateway.Route{},
		Service:             apigateway.Service{},
		Latencies: apigateway.Latencies{
			Proxy:   1,
			Gateway: 2,
			Request: 3,
		},
		ClientIP:   userIP,
		StartedAt:  12346,
		ServiceID:  serviceID,
		ConsumerID: consumerID,
	})

	filesystem := mock.FileSystemMock{}
	filesystem.On("Write", m.Anything, m.Anything).Return(nil).Twice()

	driverMock := mock.DriverMock{}
	driverMock.On("GetByService", serviceID, itemsPerPage).Return(logs).Twice()

	repo := repository.NewApiGatewayLogRepository(&driverMock)

	s, _ := service.NewApiGatewayLogParserService(repo, &filesystem)

	h := handler.NewExportMetricsByServiceHandler(s)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"", serviceID}

	err := h.HandleExportMetricsByService(context.Background())

	assert.Nil(err)
}

func TestHandleExportMetricsByService_ShouldExportCorrectMetrics(t *testing.T) {
	assert := as.New(t)

	var logs []*apigateway.Log

	consumerID := "29a5a16b-e4fa-331f-9f1c-5adea563d7de"
	serviceID := "c3e86413-648a-3552-90c3-b13491ee07d6"
	userIP := "0.0.0.0"

	itemsPerPage := 1000
	numberOfLogs := 30
	latenciesValue := 42

	for i := 0; i < numberOfLogs; i++ {
		logs = append(logs, &apigateway.Log{
			Request:             apigateway.Request{},
			UpstreamURI:         "/",
			Response:            apigateway.Response{},
			AuthenticatedEntity: apigateway.AuthenticatedEntity{ConsumerID: apigateway.Consumer{UUID: consumerID}},
			Route:               apigateway.Route{},
			Service:             apigateway.Service{},
			Latencies: apigateway.Latencies{
				Proxy:   latenciesValue,
				Gateway: latenciesValue,
				Request: latenciesValue,
			},
			ClientIP:   userIP,
			StartedAt:  12345,
			ServiceID:  serviceID,
			ConsumerID: consumerID,
		})
	}

	filesystem := mock.FileSystemMock{}

	driverMock := mock.DriverMock{}
	driverMock.On("GetByService", serviceID, itemsPerPage).Return(logs).Twice()

	repo := repository.NewApiGatewayLogRepository(&driverMock)

	s, _ := service.NewApiGatewayLogParserService(repo, &filesystem)

	h := handler.NewExportMetricsByServiceHandler(s)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"", serviceID}

	var buffer bytes.Buffer
	w := csv.NewWriter(&buffer)
	defer w.Flush()

	columns := []string{"service", "request_avg", "proxy_avg", "gateway_avg"}

	separator := ';'
	w.Comma = separator
	w.WriteAll([][]string{columns})

	requestAvg := float64(latenciesValue*numberOfLogs) / float64(numberOfLogs)
	proxyAvg := float64(latenciesValue*numberOfLogs) / float64(numberOfLogs)
	gatewayAvg := float64(latenciesValue*numberOfLogs) / float64(numberOfLogs)

	metrics := []string{
		serviceID,
		fmt.Sprintf("%.2f", requestAvg),
		fmt.Sprintf("%.2f", proxyAvg),
		fmt.Sprintf("%.2f", gatewayAvg),
	}

	w.WriteAll([][]string{metrics})
	filesystem.On("Write", m.Anything, buffer.String()).Return(nil).Once()

	err := h.HandleExportMetricsByService(context.Background())

	assert.Nil(err)
}

func TestHandleExportMetricsByService_ShouldReturnErrorOnGettingLogs(t *testing.T) {
	assert := as.New(t)

	serviceID := "c3e86413-648a-3552-90c3-b13491ee07d6"

	itemsPerPage := 1000

	filesystem := mock.FileSystemMock{}
	filesystem.On("Write", m.Anything, m.Anything).Return(nil).Twice()

	driverErr := errors.New("error on getting logs")
	driverMock := mock.DriverMock{}
	driverMock.On("GetByService", serviceID, itemsPerPage).Return(nil, driverErr).Twice()

	repo := repository.NewApiGatewayLogRepository(&driverMock)

	s, _ := service.NewApiGatewayLogParserService(repo, &filesystem)

	h := handler.NewExportMetricsByServiceHandler(s)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"", serviceID}

	err := h.HandleExportMetricsByService(context.Background())

	assert.NotNil(err)
	assert.Same(driverErr, err)
}

func TestHandleExportMetricsByService_ShouldReturnErrorOnWritingLogs(t *testing.T) {
	assert := as.New(t)

	serviceID := "c3e86413-648a-3552-90c3-b13491ee07d6"
	var logs []*apigateway.Log

	itemsPerPage := 1000

	filesystemErr := errors.New("error on writing logs")
	filesystem := mock.FileSystemMock{}
	filesystem.On("Write", m.Anything, m.Anything).Return(filesystemErr).Twice()

	driverMock := mock.DriverMock{}
	driverMock.On("GetByService", serviceID, itemsPerPage).Return(logs, nil).Twice()

	repo := repository.NewApiGatewayLogRepository(&driverMock)

	s, _ := service.NewApiGatewayLogParserService(repo, &filesystem)

	h := handler.NewExportMetricsByServiceHandler(s)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"", serviceID}

	err := h.HandleExportMetricsByService(context.Background())

	assert.NotNil(err)
	assert.Same(filesystemErr, err)
}

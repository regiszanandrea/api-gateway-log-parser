// +build integration

package test

import (
	"api-gateway-log-parser/application/handler"
	"api-gateway-log-parser/application/service"
	"api-gateway-log-parser/pkg/apigateway"
	"api-gateway-log-parser/pkg/apigateway/repository"
	mock "api-gateway-log-parser/test/mocks"
	"context"
	"errors"
	as "github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"
	"os"
	"testing"
)

func TestHandleExportByService_ShouldReturnErrorWithWrongParameters(t *testing.T) {
	assert := as.New(t)

	filesystem := mock.FileSystemMock{}
	driverMock := mock.DriverMock{}

	repo := repository.NewApiGatewayLogRepository(&driverMock)

	s, _ := service.NewApiGatewayLogParserService(repo, &filesystem)

	h := handler.NewExportByServiceHandler(s)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{}

	err := h.HandleExportByService(context.Background())

	assert.NotNil(err)
	assert.Same(err, handler.ErrServiceParameterNotFound)

	os.Args = []string{"", ""}

	err = h.HandleExportByService(context.Background())

	assert.NotNil(err)
	assert.Same(err, handler.ErrServiceParameterCouldNotBeEmpty)
}

func TestHandleExportByService_ShouldExportLogs(t *testing.T) {
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
		Latencies:           apigateway.Latencies{},
		ClientIP:            userIP,
		StartedAt:           12345,
		ServiceID:           serviceID,
		ConsumerID:          consumerID,
	})

	logs = append(logs, &apigateway.Log{
		Request:             apigateway.Request{},
		UpstreamURI:         "/",
		Response:            apigateway.Response{},
		AuthenticatedEntity: apigateway.AuthenticatedEntity{ConsumerID: apigateway.Consumer{UUID: consumerID}},
		Route:               apigateway.Route{},
		Service:             apigateway.Service{},
		Latencies:           apigateway.Latencies{},
		ClientIP:            userIP,
		StartedAt:           12346,
		ServiceID:           serviceID,
		ConsumerID:          consumerID,
	})

	filesystem := mock.FileSystemMock{}
	filesystem.On("Write", m.Anything, m.Anything).Return(nil).Twice()

	driverMock := mock.DriverMock{}
	driverMock.On("GetByService", serviceID, itemsPerPage).Return(logs).Twice()

	repo := repository.NewApiGatewayLogRepository(&driverMock)

	s, _ := service.NewApiGatewayLogParserService(repo, &filesystem)

	h := handler.NewExportByServiceHandler(s)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"", serviceID}

	err := h.HandleExportByService(context.Background())

	assert.Nil(err)
}

func TestHandleExportByService_ShouldPaginateLogs(t *testing.T) {
	assert := as.New(t)

	var logs []*apigateway.Log

	consumerID := "29a5a16b-e4fa-331f-9f1c-5adea563d7de"
	serviceID := "c3e86413-648a-3552-90c3-b13491ee07d6"
	userIP := "0.0.0.0"

	itemsPerPage := 1000

	for i := 0; i < itemsPerPage+1; i++ {
		logs = append(logs, &apigateway.Log{
			Request:             apigateway.Request{},
			UpstreamURI:         "/",
			Response:            apigateway.Response{},
			AuthenticatedEntity: apigateway.AuthenticatedEntity{ConsumerID: apigateway.Consumer{UUID: consumerID}},
			Route:               apigateway.Route{},
			Service:             apigateway.Service{},
			Latencies:           apigateway.Latencies{},
			ClientIP:            userIP,
			StartedAt:           12345,
			ServiceID:           serviceID,
			ConsumerID:          consumerID,
		})
	}

	filesystem := mock.FileSystemMock{}
	filesystem.On("Write", m.Anything, m.Anything).Return(nil).Times(3)

	driverMock := mock.DriverMock{}
	driverMock.On("GetByService", serviceID, itemsPerPage).Return(logs).Times(3)

	repo := repository.NewApiGatewayLogRepository(&driverMock)

	s, _ := service.NewApiGatewayLogParserService(repo, &filesystem)

	h := handler.NewExportByServiceHandler(s)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"", serviceID}

	err := h.HandleExportByService(context.Background())

	assert.Nil(err)
}

func TestHandleExportByService_ShouldReturnErrorOnGettingLogs(t *testing.T) {
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

	h := handler.NewExportByServiceHandler(s)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"", serviceID}

	err := h.HandleExportByService(context.Background())

	assert.NotNil(err)
	assert.Same(driverErr, err)
}

func TestHandleExportByService_ShouldReturnErrorOnWritingLogs(t *testing.T) {
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

	h := handler.NewExportByServiceHandler(s)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"", serviceID}

	err := h.HandleExportByService(context.Background())

	assert.NotNil(err)
	assert.Same(filesystemErr, err)
}

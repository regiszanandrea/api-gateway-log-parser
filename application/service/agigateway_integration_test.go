// +build integration

package service

import (
	"api-gateway-log-parser/pkg/apigateway"
	"api-gateway-log-parser/pkg/apigateway/repository"
	mock "api-gateway-log-parser/test/mocks"
	"bytes"
	"encoding/csv"
	"errors"
	as "github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"
	"strings"
	"sync"
	"testing"
)

func TestApiGatewayLogService_ShouldWriteColumns(t *testing.T) {
	assert := as.New(t)

	filesystem := mock.FileSystemMock{}

	service, _ := NewApiGatewayLogParserService(nil, &filesystem)

	var buffer bytes.Buffer
	w := csv.NewWriter(&buffer)
	defer w.Flush()

	columns := []string{
		"request",
		"upstream_uri",
		"response",
		"authenticated_entity",
		"route",
		"service",
		"latencies",
		"client_ip",
		"started_at",
		"service_id",
		"consumer_id",
	}

	columnsStr := strings.Join(columns, ";") + "\n"

	filesystem.On("Write", m.Anything, columnsStr).Return(nil).Once()

	err := service.writeColumns(w, "test.csv", &buffer)

	assert.Nil(err)
}

func TestApiGatewayLogService_ShouldReturnErrorOnWriteColumns(t *testing.T) {
	assert := as.New(t)

	filesystem := mock.FileSystemMock{}

	service, _ := NewApiGatewayLogParserService(nil, &filesystem)

	var buffer bytes.Buffer
	w := csv.NewWriter(&buffer)
	defer w.Flush()

	columns := []string{
		"request",
		"upstream_uri",
		"response",
		"authenticated_entity",
		"route",
		"service",
		"latencies",
		"client_ip",
		"started_at",
		"service_id",
		"consumer_id",
	}

	columnsStr := strings.Join(columns, ";") + "\n"

	filesystemErr := errors.New("error on writing file")

	filesystem.On("Write", m.Anything, columnsStr).Return(filesystemErr).Once()

	err := service.writeColumns(w, "test.csv", &buffer)

	assert.NotNil(err)
	assert.Same(err, filesystemErr)
}

func TestApiGatewayLogService_ShouldAddLogs(t *testing.T) {
	var logs []*apigateway.Log

	consumerID := "29a5a16b-e4fa-331f-9f1c-5adea563d7de"
	serviceID := "c3e86413-648a-3552-90c3-b13491ee07d6"
	userIP := "0.0.0.0"

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

	driverMock := mock.DriverMock{}
	repo := repository.NewApiGatewayLogRepository(&driverMock)
	driverMock.On("AddBatch", logs).Return(nil).Once()

	service, _ := NewApiGatewayLogParserService(repo, nil)

	var wg sync.WaitGroup
	wg.Add(1)

	service.addLogs(logs, &wg)

	wg.Wait()
}

func TestApiGatewayLogService_ShouldReturnErrorOnAddLogs(t *testing.T) {
	var logs []*apigateway.Log

	driverErr := errors.New("error on writing file")

	driverMock := mock.DriverMock{}
	repo := repository.NewApiGatewayLogRepository(&driverMock)
	driverMock.On("AddBatch", logs).Return(driverErr).Once()

	service, _ := NewApiGatewayLogParserService(repo, nil)

	var wg sync.WaitGroup
	wg.Add(1)

	service.addLogs(logs, &wg)

	wg.Wait()
}

func TestApiGatewayLogService_ShouldWriteLogsToFile(t *testing.T) {
	assert := as.New(t)

	var logs []*apigateway.Log

	consumerID := "29a5a16b-e4fa-331f-9f1c-5adea563d7de"
	serviceID := "c3e86413-648a-3552-90c3-b13491ee07d6"
	userIP := "0.0.0.0"

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

	fileName := "teste.csv"

	filesystem := mock.FileSystemMock{}

	var bufferTest bytes.Buffer
	wTest := csv.NewWriter(&bufferTest)
	defer wTest.Flush()

	values := getValuesFromLogs(logs)
	wTest.WriteAll(values)

	var buffer bytes.Buffer
	w := csv.NewWriter(&buffer)
	defer w.Flush()

	filesystem.On("Write", fileName, bufferTest.String()).Return(nil).Once()

	service, _ := NewApiGatewayLogParserService(nil, &filesystem)

	err := service.writeLogsToFile(logs, w, fileName, &buffer)

	assert.Nil(err)
}

func TestApiGatewayLogService_ShouldReturnErrorOnWriteLogsToFile(t *testing.T) {
	assert := as.New(t)

	var logs []*apigateway.Log
	fileName := "teste.csv"

	filesystem := mock.FileSystemMock{}

	var bufferTest bytes.Buffer
	wTest := csv.NewWriter(&bufferTest)
	defer wTest.Flush()

	values := getValuesFromLogs(logs)
	wTest.WriteAll(values)

	var buffer bytes.Buffer
	w := csv.NewWriter(&buffer)
	defer w.Flush()

	filesystemErr := errors.New("error on writing file")
	filesystem.On("Write", fileName, bufferTest.String()).Return(filesystemErr).Once()

	service, _ := NewApiGatewayLogParserService(nil, &filesystem)

	err := service.writeLogsToFile(logs, w, fileName, &buffer)

	assert.NotNil(err)
	assert.Same(err, filesystemErr)
}

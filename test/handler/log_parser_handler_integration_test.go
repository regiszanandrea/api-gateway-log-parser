// +build integration

package test

import (
	"api-gateway-log-parser/application/handler"
	"api-gateway-log-parser/application/service"
	"api-gateway-log-parser/pkg/apigateway/repository"
	mock "api-gateway-log-parser/test/mocks"
	"bufio"
	"bytes"
	"context"
	as "github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"
	"os"
	"testing"
)

func TestHandleLogParser_ShouldReturnErrorWithWrongParameters(t *testing.T) {
	assert := as.New(t)

	filesystem := mock.FileSystemMock{}
	driverMock := mock.DriverMock{}

	repo := repository.NewApiGatewayLogRepository(&driverMock)

	s, _ := service.NewApiGatewayLogParserService(repo, &filesystem)

	h := handler.NewLogParserHandler(s)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{}

	err := h.HandleApiGatewayLogParser(context.Background())

	assert.NotNil(err)
	assert.Same(err, handler.ErrPathParameterNotFound)

	os.Args = []string{"", ""}

	err = h.HandleApiGatewayLogParser(context.Background())

	assert.NotNil(err)
	assert.Same(err, handler.ErrPathParameterCouldNotBeEmpty)
}

func TestHandleLogParser_ShouldParseLogs(t *testing.T) {
	assert := as.New(t)

	buffer := bytes.NewBufferString(getLog())

	filesystem := mock.FileSystemMock{}
	driverMock := mock.DriverMock{}

	repo := repository.NewApiGatewayLogRepository(&driverMock)

	s, _ := service.NewApiGatewayLogParserService(repo, &filesystem)

	h := handler.NewLogParserHandler(s)

	path := "logs.txt"

	file := os.File{}
	scanner := bufio.NewScanner(buffer)

	filesystem.On("Open", path).Return(&file).Once()
	filesystem.On("GetScanner", &file).Return(scanner).Once()
	filesystem.On("GetLine", scanner).Return(getLog()).Twice()

	driverMock.On("AddBatch", m.Anything).Return(nil).Once()

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"", path}

	err := h.HandleApiGatewayLogParser(context.Background())

	assert.Nil(err)

}

func getLog() string {
	return `{
	  "request": {
		"method": "GET",
		"uri": "/",
		"url": "http://yost.com",
		"size": 174,
		"querystring": [],
		"headers": {
		  "accept": "*/*",
		  "host": "yost.com",
		  "user-agent": "curl/7.37.1"
		}
	  },
	  "upstream_uri": "/",
	  "response": {
		"status": 500,
		"size": 878,
		"headers": {
		  "Content-Length": "197",
		  "via": "gateway/1.3.0",
		  "Connection": "close",
		  "access-control-allow-credentials": "true",
		  "Content-Type": "application/json",
		  "server": "nginx",
		  "access-control-allow-origin": "*"
		}
	  },
	  "authenticated_entity": {
		"consumer_id": {
		  "uuid": "72b34d31-4c14-3bae-9cc6-516a0939c9d6"
		}
	  },
	  "route": {
		"created_at": 1564823899,
		"hosts": "miller.com",
		"id": "0636a119-b7ee-3828-ae83-5f7ebbb99831",
		"methods": [
		  "GET",
		  "POST",
		  "PUT",
		  "DELETE",
		  "PATCH",
		  "OPTIONS",
		  "HEAD"
		],
		"paths": [
		  "/"
		],
		"preserve_host": false,
		"protocols": [
		  "http",
		  "https"
		],
		"regex_priority": 0,
		"service": {
		  "id": "c3e86413-648a-3552-90c3-b13491ee07d6"
		},
		"strip_path": true,
		"updated_at": 1564823899
	  },
	  "service": {
		"connect_timeout": 60000,
		"created_at": 1563589483,
		"host": "ritchie.com",
		"id": "c3e86413-648a-3552-90c3-b13491ee07d6",
		"name": "ritchie",
		"path": "/",
		"port": 80,
		"protocol": "http",
		"read_timeout": 60000,
		"retries": 5,
		"updated_at": 1563589483,
		"write_timeout": 60000
	  },
	  "latencies": {
		"proxy": 1836,
		"gateway": 8,
		"request": 1058
	  },
	  "client_ip": "75.241.168.121",
	  "started_at": 1566660387
	}`
}

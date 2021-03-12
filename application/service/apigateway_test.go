package service

import (
	"api-gateway-log-parser/pkg/apigateway"
	"fmt"
	as "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestApiGatewayLogService_ShouldGetValuesFromLogs(t *testing.T) {
	assert := as.New(t)

	var logs []*apigateway.Log

	consumerID := "29a5a16b-e4fa-331f-9f1c-5adea563d7de"
	serviceID := "c3e86413-648a-3552-90c3-b13491ee07d6"
	secondServiceID := "c3e86413-648a-3552-90c3-b13491ee07d6"
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
		ServiceID:           secondServiceID,
		ConsumerID:          consumerID,
	})

	values := getValuesFromLogs(logs)

	authenticatedEntityJson := fmt.Sprintf(`{"consumer_id":{"uuid":"%s"}}`, consumerID)

	assert.IsType(values, [][]string{})
	assert.EqualValues(values[0][1], "/")
	assert.EqualValues(values[0][7], userIP)
	assert.EqualValues(values[0][9], serviceID)
	assert.EqualValues(values[0][10], consumerID)
	assert.EqualValues(values[1][9], secondServiceID)
	assert.EqualValues(values[1][10], consumerID)
	assert.EqualValues(values[1][3], authenticatedEntityJson)
}

func TestApiGatewayLogService_ShouldGetEmptyValuesFromLogs(t *testing.T) {
	assert := as.New(t)

	var logs []*apigateway.Log

	values := getValuesFromLogs(logs)

	assert.Len(values, 0)
}

func TestApiGatewayLogService_ShouldReturnFileName(t *testing.T) {
	assert := as.New(t)

	day := time.Now().Format("02-01-2006")

	serviceID := "c3e86413-648a-3552-90c3-b13491ee07d6"
	fileName := generateFileName("service", serviceID)

	assert.Contains(fileName, serviceID)
	assert.Contains(fileName, "service")
	assert.Contains(fileName, day)
}

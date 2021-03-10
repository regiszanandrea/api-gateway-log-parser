package repository

import (
	"api-gateway-log-parser/pkg/apigateway"
	"api-gateway-log-parser/pkg/apigateway/repository/driver"
)

type ApiGatewayLogRepository struct {
	driver driver.ApiGatewayLogDriver
}

func NewApiGatewayLogRepository(driver driver.ApiGatewayLogDriver) *ApiGatewayLogRepository {
	return &ApiGatewayLogRepository{
		driver: driver,
	}
}

func (a *ApiGatewayLogRepository) Add(log ...*apigateway.Log) error {
	return a.driver.AddBatch(log...)
}

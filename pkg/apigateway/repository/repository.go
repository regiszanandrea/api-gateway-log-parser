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

func (a *ApiGatewayLogRepository) GetByService(service string, limit int) ([]*apigateway.Log, error) {
	return a.driver.GetByService(service, limit)
}

func (a *ApiGatewayLogRepository) GetByConsumer(consumer string, limit int) ([]*apigateway.Log, error) {
	return a.driver.GetByConsumer(consumer, limit)
}

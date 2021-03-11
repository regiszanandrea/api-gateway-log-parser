package driver

import (
	"api-gateway-log-parser/pkg/apigateway"
)

type ApiGatewayLogDriver interface {
	GetTableName() string
	Client() interface{}
	Add(log *apigateway.Log) error
	AddBatch(...*apigateway.Log) error
	GetByService(service string, limit int) ([]*apigateway.Log, error)
}

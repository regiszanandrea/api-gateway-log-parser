package driver

import (
	"api-gateway-log-parser/pkg/apigateway"
)

type ApiGatewayLogDriver interface {
	Add(log *apigateway.Log) error
	AddBatch(...*apigateway.Log) error
	GetTableName() string
	Client() interface{}
}

package mock

import (
	"api-gateway-log-parser/pkg/apigateway"
	"github.com/stretchr/testify/mock"
)

type DriverMock struct {
	mock.Mock
}

func (d *DriverMock) GetTableName() string {
	args := d.Called()

	if args.Get(0) == nil {
		return ""
	}

	return args.Get(0).(string)
}

func (d *DriverMock) Client() interface{} {
	args := d.Called()

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0)
}

func (d *DriverMock) Add(log *apigateway.Log) error {
	args := d.Called(log)

	if args.Get(0) == nil {
		return args.Error(0)
	}

	return nil
}

func (d *DriverMock) AddBatch(logs ...*apigateway.Log) error {
	args := d.Called(logs)

	if args.Get(0) == nil {
		return args.Error(0)
	}

	return nil
}

func (d *DriverMock) GetByService(serviceID string, limit int) ([]*apigateway.Log, error) {
	args := d.Called(serviceID, limit)

	if len(d.Calls) == 2 {
		return nil, nil
	}

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*apigateway.Log), nil
}

func (d *DriverMock) GetByConsumer(consumerID string, limit int) ([]*apigateway.Log, error) {
	args := d.Called(consumerID, limit)

	if len(d.Calls) == 2 {
		return nil, nil
	}

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*apigateway.Log), nil
}

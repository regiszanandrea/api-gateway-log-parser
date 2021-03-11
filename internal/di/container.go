package di

import (
	"api-gateway-log-parser/application/handler"
	"api-gateway-log-parser/application/service"
	"api-gateway-log-parser/pkg/apigateway"
	"api-gateway-log-parser/pkg/apigateway/repository"
	"api-gateway-log-parser/pkg/apigateway/repository/driver"
	"api-gateway-log-parser/pkg/filesystem"
	"context"
	"fmt"
	"os"
)

var (
	apiGatewayLogsTableName = os.Getenv("API_GATEWAY_LOGS_TABLE_NAME_TABLE")
	dynamoURL               = os.Getenv("DYNAMODB_URL")
	dynamoRegion            = os.Getenv("DYNAMODB_REGION")
	consumerIndex           = os.Getenv("DYNAMODB_CONSUMER_INDEX")
)

type Container struct {
	logParserHandler              func(c context.Context) error
	exportByServiceHandler        func(c context.Context) error
	exportByConsumerHandler       func(c context.Context) error
	exportMetricsByServiceHandler func(c context.Context) error
	apiGatewayRepository          *repository.ApiGatewayLogRepository
	apiGatewayLogService          *service.ApiGatewayLogService
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) GetLogParserHandler() func(c context.Context) error {
	if c.logParserHandler == nil {
		c.logParserHandler = handler.NewLogParserHandler(c.MustGetApiGatewayLogService()).HandleApiGatewayLogParser
	}

	return c.logParserHandler
}

func (c *Container) GetExportByServiceHandler() func(c context.Context) error {
	if c.exportByServiceHandler == nil {
		c.exportByServiceHandler = handler.NewExportByServiceHandler(c.MustGetApiGatewayLogService()).HandleExportByService
	}

	return c.exportByServiceHandler
}

func (c *Container) GetExportMetricsByServiceHandler() func(c context.Context) error {
	if c.exportMetricsByServiceHandler == nil {
		c.exportMetricsByServiceHandler = handler.NewExportMetricsByServiceHandler(c.MustGetApiGatewayLogService()).HandleExportMetricsByService
	}

	return c.exportMetricsByServiceHandler
}

func (c *Container) GetExportByConsumerHandler() func(c context.Context) error {
	if c.exportByConsumerHandler == nil {
		c.exportByConsumerHandler = handler.NewExportByConsumerHandler(c.MustGetApiGatewayLogService()).HandleExportByConsumer
	}

	return c.exportByConsumerHandler
}

func (c *Container) MustGetApiGatewayLogService() apigateway.LogService {
	s, err := c.GetApiGatewayLogService()
	if err != nil {
		panic(err)
	}

	return s
}

func (c *Container) GetApiGatewayLogService() (*service.ApiGatewayLogService, error) {
	if c.apiGatewayLogService == nil {
		f := filesystem.NewLocalFileSystem()

		s, err := service.NewLogParserService(c.MustGetApiGatewayLogRepository(), f)
		if err != nil {
			return nil, err
		}

		c.apiGatewayLogService = s
	}

	return c.apiGatewayLogService, nil
}

func (c *Container) MustGetApiGatewayLogRepository() *repository.ApiGatewayLogRepository {
	repo, err := c.GetApiGatewayLogRepository()
	if err != nil {
		panic(err)
	}

	return repo
}

func (c *Container) GetApiGatewayLogRepository() (*repository.ApiGatewayLogRepository, error) {
	if c.apiGatewayRepository == nil {

		d, err := c.GetApiGatewayLogDriver()

		c.apiGatewayRepository = repository.NewApiGatewayLogRepository(d)

		if err != nil {
			return nil, err
		}

	}

	return c.apiGatewayRepository, nil
}

func (c *Container) GetApiGatewayLogDriver() (driver.ApiGatewayLogDriver, error) {

	fmt.Println(dynamoURL, dynamoRegion)
	d, err := driver.NewDynamoDBDriver(
		apiGatewayLogsTableName,
		driver.CreateDynamoSess(dynamoURL, dynamoRegion),
		consumerIndex,
	)

	if err != nil {
		return nil, err
	}

	return d, nil
}

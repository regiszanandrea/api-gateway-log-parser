package service

import (
	"api-gateway-log-parser/pkg/apigateway"
	"api-gateway-log-parser/pkg/apigateway/repository"
	"api-gateway-log-parser/pkg/filesystem"
	"encoding/json"
	"log"
	"sync"
)

type ApiGatewayLogService struct {
	repo       *repository.ApiGatewayLogRepository
	filesystem filesystem.API
}

func NewLogParserService(repo *repository.ApiGatewayLogRepository, filesystem filesystem.API) (*ApiGatewayLogService, error) {
	return &ApiGatewayLogService{
		repo:       repo,
		filesystem: filesystem,
	}, nil
}

func (a *ApiGatewayLogService) Parse(path string) error {
	file, err := a.filesystem.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := a.filesystem.GetScanner(file)

	var logs []*apigateway.Log

	var wg sync.WaitGroup

	i := 0
	logsBatchMaxLen := 200

	for scanner.Scan() {
		if i > 200 {
			break
		}
		var apiGatewayLog apigateway.Log

		line := []byte(a.filesystem.GetLine(scanner))

		err = json.Unmarshal(line, &apiGatewayLog)

		if err != nil {
			return err
		}

		apiGatewayLog.ServiceID = apiGatewayLog.Service.ID
		apiGatewayLog.CustomerID = apiGatewayLog.AuthenticatedEntity.ConsumerID.UUID

		logs = append(logs, &apiGatewayLog)

		if len(logs) > logsBatchMaxLen {
			wg.Add(1)

			a.addLogs(logs, &wg)

			logs = nil
		}
		i += 1
	}

	wg.Wait()

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (a *ApiGatewayLogService) addLogs(logs []*apigateway.Log, wg *sync.WaitGroup) {
	err := a.repo.Add(logs...)
	if err != nil {
		log.Fatal(err)
	}

	defer wg.Done()
}

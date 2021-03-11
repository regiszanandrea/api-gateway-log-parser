package service

import (
	"api-gateway-log-parser/pkg/apigateway"
	"api-gateway-log-parser/pkg/apigateway/repository"
	"api-gateway-log-parser/pkg/filesystem"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

const itemsPerPage = 1000

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

	logsBatchMaxLen := 200

	for scanner.Scan() {
		var apiGatewayLog apigateway.Log

		line := []byte(a.filesystem.GetLine(scanner))

		err = json.Unmarshal(line, &apiGatewayLog)

		if err != nil {
			return err
		}

		apiGatewayLog.ServiceID = apiGatewayLog.Service.ID
		apiGatewayLog.ConsumerID = apiGatewayLog.AuthenticatedEntity.ConsumerID.UUID

		logs = append(logs, &apiGatewayLog)

		if len(logs) > logsBatchMaxLen {
			wg.Add(1)

			a.addLogs(logs, &wg)

			logs = nil
		}
	}

	wg.Wait()

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (a *ApiGatewayLogService) ExportByService(service string) error {
	fileName := generateFileName("service", service)

	var buffer bytes.Buffer
	w := csv.NewWriter(&buffer)
	defer w.Flush()

	err := a.writeColumns(w, fileName, &buffer)
	if err != nil {
		return err
	}

	for {
		logs, err := a.repo.GetByService(service, itemsPerPage)

		if err != nil {
			return err
		}

		if logs == nil {
			break
		}

		err = a.WriteLogsToFile(logs, w, fileName, &buffer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *ApiGatewayLogService) ExportByConsumer(consumer string) error {
	fileName := generateFileName("consumer", consumer)

	var buffer bytes.Buffer

	w := csv.NewWriter(&buffer)
	defer w.Flush()

	err := a.writeColumns(w, fileName, &buffer)
	if err != nil {
		return err
	}

	for {
		logs, err := a.repo.GetByConsumer(consumer, itemsPerPage)

		if err != nil {
			return err
		}

		if logs == nil {
			break
		}

		err = a.WriteLogsToFile(logs, w, fileName, &buffer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *ApiGatewayLogService) WriteLogsToFile(logs []*apigateway.Log, w *csv.Writer, fileName string, buffer *bytes.Buffer) error {
	values := getValuesFromLogs(logs)

	err := w.WriteAll(values)

	if err != nil {
		return err
	}

	err = a.filesystem.Write(fileName, buffer.String())

	if err != nil {
		return err
	}

	buffer.Reset()
	return nil
}

func (a *ApiGatewayLogService) addLogs(logs []*apigateway.Log, wg *sync.WaitGroup) {
	err := a.repo.Add(logs...)
	if err != nil {
		log.Fatal(err)
	}

	defer wg.Done()
}

func generateFileName(prefix string, id string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%s-%s-%s-%d.csv", prefix, id, time.Now().Format("02-01-2006"), rand.Uint32())
}

func (a *ApiGatewayLogService) writeColumns(w *csv.Writer, fileName string, buffer *bytes.Buffer) error {
	columns := apigateway.GetJsonFieldsFromLogStruct()
	separator := ';'
	w.Comma = separator
	err := w.WriteAll([][]string{columns})

	if err != nil {
		return err
	}

	err = a.filesystem.Write(fileName, buffer.String())

	if err != nil {
		return err
	}

	defer w.Flush()
	defer buffer.Reset()

	return nil
}

func getValuesFromLogs(logs []*apigateway.Log) [][]string {
	var values [][]string
	for _, l := range logs {
		r := l.ToSlice()

		values = append(values, r)
	}
	return values
}

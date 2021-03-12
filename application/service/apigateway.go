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

func NewApiGatewayLogParserService(repo *repository.ApiGatewayLogRepository, filesystem filesystem.API) (*ApiGatewayLogService, error) {
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

		if len(line) == 0 {
			break
		}

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

		err = a.writeLogsToFile(logs, w, fileName, &buffer)
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

		err = a.writeLogsToFile(logs, w, fileName, &buffer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *ApiGatewayLogService) ExportMetricsByService(service string) error {
	fileName := generateFileName("metrics", service)

	var buffer bytes.Buffer
	w := csv.NewWriter(&buffer)
	defer w.Flush()

	columns := []string{"service", "request_avg", "proxy_avg", "gateway_avg"}

	separator := ';'
	w.Comma = separator
	err := w.WriteAll([][]string{columns})

	if err != nil {
		return err
	}

	var requestAvg, proxyAvg, gatewayAvg float64
	var requestSum, proxySum, gatewaySum int

	numberOfLogs := 0

	for {
		logs, err := a.repo.GetByService(service, itemsPerPage)

		if err != nil {
			return err
		}

		if logs == nil {
			break
		}

		for _, l := range logs {
			requestSum += l.Latencies.Request
			proxySum += l.Latencies.Proxy
			gatewaySum += l.Latencies.Gateway

			numberOfLogs += 1
		}
	}

	requestAvg = float64(requestSum) / float64(numberOfLogs)
	proxyAvg = float64(proxySum) / float64(numberOfLogs)
	gatewayAvg = float64(gatewaySum) / float64(numberOfLogs)

	metrics := []string{
		service,
		fmt.Sprintf("%.2f", requestAvg),
		fmt.Sprintf("%.2f", proxyAvg),
		fmt.Sprintf("%.2f", gatewayAvg),
	}

	err = w.WriteAll([][]string{metrics})

	if err != nil {
		return err
	}

	err = a.filesystem.Write(fileName, buffer.String())

	if err != nil {
		return err
	}

	return nil
}

func (a *ApiGatewayLogService) writeLogsToFile(logs []*apigateway.Log, w *csv.Writer, fileName string, buffer *bytes.Buffer) error {
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

func generateFileName(prefix string, id string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%s-%s-%s-%d.csv", prefix, id, time.Now().Format("02-01-2006"), rand.Uint32())
}

func getValuesFromLogs(logs []*apigateway.Log) [][]string {
	var values [][]string
	for _, l := range logs {
		r := l.ToSlice()

		values = append(values, r)
	}
	return values
}

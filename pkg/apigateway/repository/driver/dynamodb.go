package driver

import (
	"api-gateway-log-parser/pkg/apigateway"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strconv"
)

func CreateDynamoSess(url string, region string) *dynamodb.DynamoDB {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess, &aws.Config{
		Endpoint: aws.String(url),
		Region:   aws.String(region),
	})
}

type dynamoDB struct {
	db               *dynamodb.DynamoDB
	tableName        string
	consumerIndex    string
	startKey         map[string]*dynamodb.AttributeValue
	lastPageAchieved bool
}

func NewDynamoDBDriver(tableName string, db *dynamodb.DynamoDB, consumerIndex string) (ApiGatewayLogDriver, error) {
	return &dynamoDB{
		tableName:     tableName,
		db:            db,
		consumerIndex: consumerIndex,
	}, nil
}

func (d *dynamoDB) Add(log *apigateway.Log) error {
	panic("implement me")
}

func (d *dynamoDB) AddBatch(logs ...*apigateway.Log) error {
	batchSize := 25

	logsCount := len(logs)

	for i := 0; i < logsCount; i += batchSize {
		low := i
		high := low + batchSize

		if high > logsCount {
			high = logsCount
		}

		l := logs[low:high]
		var writeRequests []*dynamodb.WriteRequest

		for _, log := range l {
			item, err := dynamodbattribute.MarshalMap(&log)
			if err != nil {
				return err
			}
			writeRequests = append(writeRequests, &dynamodb.WriteRequest{
				PutRequest: &dynamodb.PutRequest{Item: item},
			})
		}

		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				d.tableName: writeRequests,
			},
		}

		_, err := d.db.BatchWriteItem(input)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *dynamoDB) Client() interface{} {
	return d.db
}

func (d *dynamoDB) GetTableName() string {
	return d.tableName
}

func (d *dynamoDB) GetByService(serviceID string, limit int) ([]*apigateway.Log, error) {
	if d.lastPageAchieved {
		return nil, nil
	}

	input := &dynamodb.QueryInput{
		TableName: &d.tableName,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":value": {
				S: &serviceID,
			},
		},
		Limit:                  aws.Int64(int64(limit)),
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :value", "service_id")),
	}

	return d.getLogsByQuery(input)
}

func (d *dynamoDB) GetByConsumer(consumerID string, limit int) ([]*apigateway.Log, error) {
	if d.lastPageAchieved {
		return nil, nil
	}

	input := &dynamodb.QueryInput{
		TableName: &d.tableName,
		IndexName: aws.String(d.consumerIndex),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":value": {
				S: &consumerID,
			},
		},
		Limit:                  aws.Int64(int64(limit)),
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :value", "consumer_id")),
	}

	return d.getLogsByQuery(input)
}

func (d *dynamoDB) getLogsByQuery(input *dynamodb.QueryInput) ([]*apigateway.Log, error) {
	if d.startKey != nil {
		input.ExclusiveStartKey = d.startKey
	}

	var logs []*apigateway.Log

	result, err := d.db.Query(input)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return d.convertToLogs(result, logs)
}

func (d *dynamoDB) convertToLogs(result *dynamodb.QueryOutput, logs []*apigateway.Log) ([]*apigateway.Log, error) {
	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &logs)

	if err != nil {
		return nil, err
	}

	if result.LastEvaluatedKey == nil {
		d.lastPageAchieved = true
		return logs, nil
	}

	var lastLog *apigateway.Log
	err = dynamodbattribute.UnmarshalMap(result.LastEvaluatedKey, &lastLog)
	if err != nil {
		return nil, err
	}

	startedAt := strconv.Itoa(int(lastLog.StartedAt))

	if lastLog.ConsumerID != "" {
		d.startKey = generateConsumerStartKey(lastLog, startedAt)
		return logs, nil
	}

	d.startKey = generateServiceStartKey(lastLog, startedAt)
	return logs, nil
}

func generateServiceStartKey(lastLog *apigateway.Log, startedAt string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"service_id": {
			S: aws.String(lastLog.ServiceID),
		},
		"started_at": {
			N: aws.String(startedAt),
		},
	}
}

func generateConsumerStartKey(lastLog *apigateway.Log, startedAt string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"service_id": {
			S: aws.String(lastLog.ServiceID),
		},
		"started_at": {
			N: aws.String(startedAt),
		},
		"consumer_id": {
			S: aws.String(lastLog.ConsumerID),
		},
	}
}

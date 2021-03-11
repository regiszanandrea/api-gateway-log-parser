package driver

import (
	"api-gateway-log-parser/pkg/apigateway"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
	"strconv"
)

func CreateDynamoSess(url string, region string) *dynamodb.DynamoDB {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess, &aws.Config{
		Endpoint: aws.String(os.Getenv("DYNAMODB_URL")),
		Region:   aws.String(os.Getenv("DYNAMODB_REGION")),
	})
}

type dynamoDB struct {
	tableName        string
	db               *dynamodb.DynamoDB
	startKey         map[string]*dynamodb.AttributeValue
	lastPageAchieved bool
}

func NewDynamoDBDriver(tableName string, db *dynamodb.DynamoDB) (ApiGatewayLogDriver, error) {
	return &dynamoDB{
		tableName: tableName,
		db:        db,
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

func (d *dynamoDB) GetByService(service string, limit int) ([]*apigateway.Log, error) {
	if d.lastPageAchieved {
		return nil, nil
	}

	input := &dynamodb.QueryInput{
		TableName: &d.tableName,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":value": {
				S: &service,
			},
		},
		Limit:                  aws.Int64(int64(limit)),
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :value", "service_id")),
	}

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

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &logs)

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

	d.startKey = generateStartKey(lastLog, startedAt)

	return logs, nil
}

func generateStartKey(lastLog *apigateway.Log, startedAt string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"service_id": {
			S: aws.String(lastLog.ServiceID),
		},
		"started_at": {
			N: aws.String(startedAt),
		},
	}
}

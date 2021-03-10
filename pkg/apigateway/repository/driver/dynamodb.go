package driver

import (
	"api-gateway-log-parser/pkg/apigateway"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
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
	tableName string
	db        *dynamodb.DynamoDB
}

func NewDynamoDBDriver(tableName string, db *dynamodb.DynamoDB) (ApiGatewayLogDriver, error) {
	return &dynamoDB{
		tableName,
		db,
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

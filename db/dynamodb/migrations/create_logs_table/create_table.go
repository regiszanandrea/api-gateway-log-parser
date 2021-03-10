package main

import (
	"api-gateway-log-parser/internal/di"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
)

var container *di.Container

func init() {
	container = di.NewContainer()
}

func main() {
	d, err := container.GetApiGatewayLogDriver()

	if err != nil {
		log.Fatal(err)
	}

	dbSvc := d.Client().(*dynamodb.DynamoDB)

	log.Println("creating table")
	params := &dynamodb.CreateTableInput{
		TableName: aws.String(d.GetTableName()),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("service_id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("consumer_id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("started_at"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("service_id"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("started_at"),
				KeyType:       aws.String("RANGE"),
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("ConsumerIDIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("consumer_id"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String(dynamodb.ProjectionTypeAll),
				},
			},
		},
	}

	_, err = dbSvc.CreateTable(params)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("table created successfully")
}

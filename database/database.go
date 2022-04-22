package database

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"fmt"
)

type User struct {
	Username string `json:"username"`
	PassHash string `json:"passhash"`
}

const tableName = "goRateUsers"

func initDynamo() *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("/home/tristan/.aws/credentials", ""),
	})
	if err != nil {
		fmt.Printf("%s", err)
	}
	svc := dynamodb.New(sess)
	return svc
}

func AddUser(u User) {
	svc := initDynamo()
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		log.Fatalf("Got error marshalling new user item: %s", err)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
}

func GetUser(username string) User {
	svc := initDynamo()
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		log.Fatalf("Could not retrieve skill")
	}
	item := User{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	return item
}

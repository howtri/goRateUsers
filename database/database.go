package database

import (
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"fmt"
)

type Skill struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Rankings []int  `json:"rankings2"`
	Ranking  int    `json:"ranking"`
}

const tableName = "goRateSkills"

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

func GetTable() {
	svc := initDynamo()
	req := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
	result, err := svc.DescribeTable(req)
	if err != nil {
		fmt.Printf("%s", err)
	}
	table := result.Table
	fmt.Printf("done", table)
}

func AddSkill(s Skill) {
	svc := initDynamo()
	s2 := []int{0}
	data := Skill{ID: s.ID, Name: s.Name, Rankings: s2}
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		log.Fatalf("Got error marshalling new skill item: %s", err)
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

func GetSkill(id string) Skill {
	svc := initDynamo()
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		log.Fatalf("Could not retrieve skill")
	}
	item := Skill{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	return item
}

func SearchSkills(s Skill) []Skill {
	svc := initDynamo()
	proj := expression.NamesList(expression.Name("id"), expression.Name("name"), expression.Name("rankings"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	result, err := svc.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}

	var results []Skill

	for _, v := range result.Items {
		item := Skill{}
		err = dynamodbattribute.UnmarshalMap(v, &item)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}

		if strings.Contains(item.Name, s.Name) {
			results = append(results, item)
		}
	}
	return results
}

func RankSkill(s Skill) {
	svc := initDynamo()
	av := &dynamodb.AttributeValue{
		N: aws.String(strconv.Itoa(s.Ranking)),
	}

	var qids []*dynamodb.AttributeValue
	qids = append(qids, av)
	log.Printf("For the id %s", s.ID)
	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(s.ID),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":qid": {
				L: qids,
			},
			":empty_list": {
				L: []*dynamodb.AttributeValue{},
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET rankings2 = list_append(if_not_exists(rankings2, :empty_list), :qid)"),
		TableName:        aws.String(tableName),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		log.Fatalf("Got error calling UpdateItem: %s", err)
	}
}

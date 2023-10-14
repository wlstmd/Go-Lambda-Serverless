package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type MyEvent struct {
	Subject   string `json:"subject"`
	Image_url string `json:"image_url"`
}

type Picture struct {
	Id        int64  `json:"id"`
	Subject   string `json:"subject"`
	Image_url string `json:"image_url"`
}

func HandleRequest(ctx context.Context, myInput MyEvent) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	picture := Picture{
		Id:        time.Now().UnixNano() / 1000000,
		Subject:   myInput.Subject,
		Image_url: myInput.Image_url,
	}

	av, err := dynamodbattribute.MarshalMap(picture)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}

	// Create item in table Movies
	tableName := "Gallery"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	return fmt.Sprintln("Successfully added '" + picture.Subject + " to table " + tableName), nil
}

func main() {
	lambda.Start(HandleRequest)
}

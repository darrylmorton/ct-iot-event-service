package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/darrylmorton/ct-iot-event-service/internal/app"
	"github.com/darrylmorton/ct-iot-event-service/internal/data"
	"log"
	"os"
)

func main() {
	var envConfig app.EnvConfig

	flag.IntVar(&envConfig.Port, "port", 4000, "API server port")
	flag.StringVar(&envConfig.Env, "env", "dev", "Environment (dev|stage|prod)")
	flag.StringVar(&envConfig.Dsn, "db-dsn", os.Getenv("EVENTS_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	dbClient, err := sql.Open("postgres", envConfig.Dsn)
	if err != nil {
		logger.Fatal(err)
	}

	queue := app.QueueName
	timeout := 20

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	sqsClient := sqs.NewFromConfig(cfg)

	gQInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queue),
	}

	urlResult, err := app.GetQueueURL(context.TODO(), sqsClient, gQInput)
	if err != nil {
		logger.Printf("Error getting queue url:%v\n", err)
	}

	queueURL := urlResult.QueueUrl

	gMInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: 25,
		VisibilityTimeout:   int32(timeout),
	}

	serviceConfig := app.ServiceConfig{
		SqsClient:              sqsClient,
		SqsReceiveMessageInput: gMInput,
		DbClient:               dbClient,
		Logger:                 logger,
		Models:                 data.NewModels(dbClient),
		EnvConfig:              envConfig,
	}

	srv := app.StartServer(&serviceConfig)
	err = srv.ListenAndServe()
	fmt.Printf("The server has encountered an error: %v", err)
}

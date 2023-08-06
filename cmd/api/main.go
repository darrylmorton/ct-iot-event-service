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

	// ** SQS START
	queue := app.QueueName
	timeout := 20

	//queue := flag.String("q", "", "The name of the queue")
	//timeout := flag.Int("t", 5, "How long, in seconds, that the message is hidden from others")
	//flag.Parse()

	//if *queue == "" {
	//	fmt.Println("You must supply the name of a queue (-q QUEUE)")
	//	//return
	//}
	//
	//if *timeout < 0 {
	//	*timeout = 0
	//}
	//
	//if *timeout > 12*60*60 {
	//	*timeout = 12 * 60 * 60
	//}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	sqsClient := sqs.NewFromConfig(cfg)

	gQInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queue),
	}

	// Get URL of queue
	urlResult, err := app.GetQueueURL(context.TODO(), sqsClient, gQInput)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		//return
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

	//logger.Println("cfg", cfg)

	serviceConfig := app.ServiceConfig{
		SqsClient:              sqsClient,
		SqsReceiveMessageInput: gMInput,
		DbClient:               dbClient,
		Logger:                 logger,
		Models:                 data.NewModels(dbClient),
		EnvConfig:              envConfig,
	}
	// ** SQS END

	srv := app.StartServer(&serviceConfig)
	err = srv.ListenAndServe()
	fmt.Printf("The server has encountered an error: %v", err)
}

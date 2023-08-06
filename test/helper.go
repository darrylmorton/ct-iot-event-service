package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/darrylmorton/ct-iot-event-service/client"
	"github.com/darrylmorton/ct-iot-event-service/internal/app"
	"github.com/darrylmorton/ct-iot-event-service/internal/data"
	"github.com/darrylmorton/ct-iot-event-service/internal/models"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

var envs = client.CreateEnvs("", "", "", "v1")

var Events = []models.Event{
	{
		DeviceId:    "esp32-2424242424",
		Description: "Maximum threshold reached",
		Type:        "temperature",
		Event:       "TEMPERATURE_MAX_THRESHOLD",
		Read:        false,
	},
	{
		DeviceId:    "esp32-4646464646",
		Description: "Minimum threshold reached",
		Type:        "temperature",
		Event:       "TEMPERATURE_MIN_THRESHOLD",
		Read:        false,
	}, {
		DeviceId:    "esp32-0123456789",
		Description: "Maximum threshold reached",
		Type:        "temperature",
		Event:       "TEMPERATURE_MAX_THRESHOLD",
		Read:        false,
	},
}

type DbConfig struct {
	Client *sql.DB
}

type SQSReceiveMessageImpl struct{}

func (dt SQSReceiveMessageImpl) GetQueueUrl(ctx context.Context,
	params *sqs.GetQueueUrlInput,
	optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error) {

	output := &sqs.GetQueueUrlOutput{
		QueueUrl: aws.String(app.QueueName),
	}

	return output, nil
}

func (dt SQSReceiveMessageImpl) ReceiveMessage(ctx context.Context,
	params *sqs.ReceiveMessageInput,
	optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {

	messageOne := Events[0]
	messageOneMarshalled, _ := json.Marshal(messageOne)

	messageTwo := Events[1]
	messageTwoMarshalled, _ := json.Marshal(messageTwo)

	messages := []types.Message{
		{
			MessageId:     aws.String("message-one-id"),
			ReceiptHandle: aws.String("message-one-receipt-handle"),
			Body:          aws.String(string(messageOneMarshalled)),
		},
		{
			MessageId:     aws.String("message-two-id"),
			ReceiptHandle: aws.String("message-two-receipt-handle"),
			Body:          aws.String(string(messageTwoMarshalled)),
		},
	}

	output := &sqs.ReceiveMessageOutput{
		Messages: messages,
	}

	return output, nil
}

func createHeaders() map[string]string {
	headers := make(map[string]string)

	headers["Accept"] = "application/json"

	return headers
}

func DbClient() *sql.DB {
	dbDsn := os.Getenv("EVENTS_DB_DSN")

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := sql.Open("postgres", dbDsn)
	if err != nil {
		logger.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("database connection pool established")

	return db
}

func (config *DbConfig) CreateDbTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS events(
			id uuid PRIMARY KEY UNIQUE DEFAULT gen_random_uuid() NOT NULL,
			device_id VARCHAR NOT NULL,
			description VARCHAR NOT NULL,
			type VARCHAR NOT NULL ,
			event VARCHAR NOT NULL,
			read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX updated_at ON events (updated_at);
	`

	_, err := config.Client.Query(query)
	if err != nil {
		fmt.Errorf("DeleteEventsErr %v", err)
		return err
	}

	return nil
}

func (config *DbConfig) DropDbTable() error {
	query := `
		DROP TABLE IF EXISTS events
	`

	_, err := config.Client.Query(query)
	if err != nil {
		fmt.Errorf("DeleteEventsErr %v", err)
		return err
	}

	return nil
}

func (config *DbConfig) DeleteEvents() error {
	query := `
		DELETE FROM events
	`

	_, err := config.Client.Query(query)
	if err != nil {
		fmt.Errorf("DeleteEventsErr %v", err)
		return err
	}

	return nil
}

func (dbConfig *DbConfig) CreateEvent(data models.Event) (models.Event, error) {
	query := `
		INSERT INTO events (device_id, description, type, event, read)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, device_id AS deviceId, description, type, event, read
	`

	var event models.Event

	args := []interface{}{data.DeviceId, data.Description, data.Type, data.Event, data.Read}
	row := dbConfig.Client.QueryRow(query, args...)

	err := row.Scan(
		&event.Id,
		&event.DeviceId,
		&event.Description,
		&event.Type,
		&event.Event,
		&event.Read,
	)

	if err != nil {
		return models.Event{}, err
	}

	return event, err
}

func StartServer() *http.Server {
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

	sqsClient := &SQSReceiveMessageImpl{}
	timeout := 20
	queueUrlInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String(app.QueueName),
	}

	urlResult, err := app.GetQueueURL(context.Background(), sqsClient, queueUrlInput)
	if err != nil {
		logger.Fatal(err)
	}
	queueURL := urlResult.QueueUrl

	receiveMessageInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: 25,
		VisibilityTimeout:   int32(timeout),
	}

	serviceConfig := app.ServiceConfig{
		SqsClient:              sqsClient,
		SqsReceiveMessageInput: receiveMessageInput,
		DbClient:               dbClient,
		Logger:                 logger,
		Models:                 data.NewModels(dbClient),
		EnvConfig:              envConfig,
	}

	return app.StartServer(&serviceConfig)
}

type HealthCheck struct {
	Version     string `json:"version,omitempty"`
	Status      string `json:"status,omitempty"`
	Environment string `json:"environment,omitempty"`
}

func GetHealthCheck() (int, HealthCheck) {
	url := fmt.Sprintf("%s/health", envs.ClientUrl)

	requestOptions := client.RequestOptions{
		Headers: createHeaders(),
		Method:  "GET",
		Url:     url,
		Payload: nil,
	}

	res := client.GetHealthCheckRequest(requestOptions)

	return GetHealthCheckResponse(res)
}

func PutEvent(id string, payload models.Event) (int, models.Event) {
	url := fmt.Sprintf("%s/events/%s", envs.ClientUrl, id)

	payloadMarshalled, _ := json.Marshal(payload)

	requestOptions := client.RequestOptions{
		Headers: createHeaders(),
		Method:  "PUT",
		Url:     url,
		Payload: payloadMarshalled,
	}

	res := client.PutRequest(requestOptions)

	return PutEventResponse(res)
}

func GetEvents() (int, []models.Event) {
	url := fmt.Sprintf("%s/events", envs.ClientUrl)

	requestOptions := client.RequestOptions{
		Headers: createHeaders(),
		Method:  "GET",
		Url:     url,
		Payload: nil,
	}

	res := client.GetEventsRequest(requestOptions)

	return GetEventsResponse(res)
}

func GetEvent(id string) (int, models.Event) {
	url := fmt.Sprintf("%s/events/%s", envs.ClientUrl, id)

	requestOptions := client.RequestOptions{
		Headers: createHeaders(),
		Method:  "GET",
		Url:     url,
		Payload: nil,
	}

	res := client.GetEventRequest(requestOptions)

	return GetEventResponse(res)
}

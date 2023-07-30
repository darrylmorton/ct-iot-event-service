package test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/darrylmorton/ct-iot-event-service/client"
	"github.com/darrylmorton/ct-iot-event-service/internal/app"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

var envs = client.CreateEnvs("", "", "", "v1")

var Events = []app.Event{
	{
		DeviceName:  "esp32-2424242424",
		Description: "Maximum threshold reached",
		Type:        "temperature",
		Event:       "TEMPERATURE_MAX_THRESHOLD",
		Read:        false,
	},
	{
		DeviceName:  "esp32-4646464646",
		Description: "Minimum threshold reached",
		Type:        "temperature",
		Event:       "TEMPERATURE_MIN_THRESHOLD",
		Read:        false,
	}, {
		DeviceName:  "esp32-0123456789",
		Description: "Maximum threshold reached",
		Type:        "temperature",
		Event:       "TEMPERATURE_MAX_THRESHOLD",
		Read:        false,
	},
}

func createHeaders() map[string]string {
	headers := make(map[string]string)

	headers["Accept"] = "application/json"

	return headers
}

func DbConnection() *sql.DB {
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

func CreateDbTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS events(
			id uuid PRIMARY KEY UNIQUE DEFAULT gen_random_uuid() NOT NULL,
			device_name VARCHAR NOT NULL,
			description VARCHAR NOT NULL,
			type VARCHAR NOT NULL ,
			event VARCHAR NOT NULL,
			read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX updated_at ON events (updated_at);
	`

	_, err := db.Query(query)
	if err != nil {
		fmt.Errorf("DeleteEventsErr %v", err)
		return err
	}

	return nil
}

func DropDbTable(db *sql.DB) error {
	query := `
		DROP TABLE IF EXISTS events
	`

	_, err := db.Query(query)
	if err != nil {
		fmt.Errorf("DeleteEventsErr %v", err)
		return err
	}

	return nil
}

func DeleteEvents(db *sql.DB) error {
	query := `
		DELETE FROM events
	`

	_, err := db.Query(query)
	if err != nil {
		fmt.Errorf("DeleteEventsErr %v", err)
		return err
	}

	return nil
}

func CreateEvent(db *sql.DB, data app.Event) (app.Event, error) {
	query := `
		INSERT INTO events (device_name, description, type, event, read)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, device_name AS deviceName, description, type, event, read
	`

	var event app.Event

	args := []interface{}{data.DeviceName, data.Description, data.Type, data.Event, data.Read}
	row := db.QueryRow(query, args...)

	err := row.Scan(
		&event.Id,
		&event.DeviceName,
		&event.Description,
		&event.Type,
		&event.Event,
		&event.Read,
	)

	if err != nil {
		return app.Event{}, err
	}

	return event, err
}

type mockedReceiveMsgs struct {
	sqsiface.SQSAPI
	Resp sqs.ReceiveMessageOutput
}

func (m mockedReceiveMsgs) ReceiveMessage(in *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	// Only need to return mocked response output
	return &m.Resp, nil
}

func CreateMessages() sqs.ReceiveMessageOutput {
	messageOne := Events[0]
	messageOneMarshalled, _ := json.Marshal(messageOne)

	messageTwo := Events[1]
	messageTwoMarshalled, _ := json.Marshal(messageTwo)

	return sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{
			{Body: aws.String(string(messageOneMarshalled))},
			{Body: aws.String(string(messageTwoMarshalled))},
		},
	}
}

func StartServer() *http.Server {
	q := app.Queue{
		Client: mockedReceiveMsgs{Resp: CreateMessages()},
		URL:    fmt.Sprintf("mockURL_%d", 0),
	}

	return app.StartServer(q)
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

func PutEvent(id string, payload app.Event) (int, app.Event) {
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

func GetEvents() (int, []app.Event) {
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

func GetEvent(id string) (int, app.Event) {
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

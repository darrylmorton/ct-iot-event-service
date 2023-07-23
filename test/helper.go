package test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/darrylmorton/ct-iot-event-service/client"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

var envs = client.CreateEnvs("", "", "", "v1/events")

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
		CREATE TABLE events(
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
		DROP TABLE events
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

func CreateEvent(db *sql.DB, data Event) (Event, error) {
	query := `
		INSERT INTO events (device_name, description, type, event, read)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, device_name AS deviceName, description, type, event, read
	`

	var event Event

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
		return Event{}, err
	}

	return event, err
}

type HealthCheck struct {
	Version     string `json:"version,omitempty"`
	Status      string `json:"status,omitempty"`
	Environment string `json:"environment,omitempty"`
}

type Event struct {
	Id          string    `json:"id,omitempty"`
	DeviceName  string    `json:"deviceName,omitempty"`
	Description string    `json:"description,omitempty"`
	Type        string    `json:"type,omitempty"`
	Event       string    `json:"event,omitempty"`
	Read        bool      `json:"read,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

func CreateEventPayload() Event {
	return Event{
		DeviceName:  "esp32-0123456789",
		Description: "Maximum threshold reached",
		Type:        "temperature",
		Event:       "TEMPERATURE_MAX_THRESHOLD",
		Read:        false,
	}
}

func GetHealthCheck() (int, HealthCheck) {
	url := fmt.Sprintf("%s/%s", envs.ClientUrl, "health")

	requestOptions := client.RequestOptions{
		Headers: createHeaders(),
		Method:  "GET",
		Url:     url,
		Payload: nil,
	}

	res := client.GetHealthCheckRequest(requestOptions)

	return GetHealthCheckResponse(res)
}

func PutEvent(id string, payload Event) (int, Event) {
	url := fmt.Sprintf("%s/%s", envs.ClientUrl, id)

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

func GetEvents() (int, []Event) {
	requestOptions := client.RequestOptions{
		Headers: createHeaders(),
		Method:  "GET",
		Url:     envs.ClientUrl,
		Payload: nil,
	}

	res := client.GetAllRequest(requestOptions)

	return GetEventsResponse(res)
}

func GetEvent(id string) (int, Event) {
	url := fmt.Sprintf("%s/%s", envs.ClientUrl, id)

	requestOptions := client.RequestOptions{
		Headers: createHeaders(),
		Method:  "GET",
		Url:     url,
		Payload: nil,
	}

	res := client.GetRequest(requestOptions)

	return GetEventResponse(res)
}

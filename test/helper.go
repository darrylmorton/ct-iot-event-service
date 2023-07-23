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

func DeleteEvents(db *sql.DB) error {
	query := `
		DELETE FROM events
	`

	_, err := db.Query(query)
	if err != nil {
		fmt.Errorf("DeleteEventsErr %v", err)
		return err
	}

	defer db.Close()

	return nil
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

func PostEvent() (int, Event) {
	payload := CreateEventPayload()
	payloadMarshalled, payloadMarshalledErr := json.Marshal(payload)
	if payloadMarshalledErr != nil {
		err := fmt.Errorf("PostEvent - payloadMarshalledErr %v", payloadMarshalledErr)
		fmt.Errorf(err.Error())
	}

	requestOptions := client.RequestOptions{
		Headers: createHeaders(),
		Method:  "POST",
		Url:     envs.ClientUrl,
		Payload: payloadMarshalled,
	}

	res := client.CreateRequest(requestOptions)

	return PostEventResponse(res)
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

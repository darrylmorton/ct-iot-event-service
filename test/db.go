package test

import (
	"database/sql"
	"fmt"
	"github.com/darrylmorton/ct-iot-event-service/internal/models"
	"log"
	"os"
)

type DbConfig struct {
	Client *sql.DB
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
		fmt.Printf("DeleteEventsErr %v", err)
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
		fmt.Printf("DeleteEventsErr %v", err)
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
		fmt.Printf("DeleteEventsErr %v", err)
		return err
	}

	return nil
}

func (config *DbConfig) CreateEvent(data models.Event) (models.Event, error) {
	query := `
		INSERT INTO events (device_id, description, type, event, read)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, device_id AS deviceId, description, type, event, read
	`

	var event models.Event

	args := []interface{}{data.DeviceId, data.Description, data.Type, data.Event, data.Read}
	row := config.Client.QueryRow(query, args...)

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

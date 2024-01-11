package data

import (
	"database/sql"
	"errors"
	"github.com/darrylmorton/ct-iot-event-service/internal/models"
	"log"
	"net/http"
)

type EventModel struct {
	DB *sql.DB
}

func (e EventModel) GetEvents() ([]*models.Event, error) {
	query := `
	  SELECT id, device_id, description, type, event, read FROM events
	  ORDER BY updated_at
  	`

	rows, err := e.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []*models.Event

	for rows.Next() {
		var event models.Event

		err := rows.Scan(
			&event.Id,
			&event.DeviceId,
			&event.Description,
			&event.Type,
			&event.Event,
			&event.Read,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (e EventModel) GetEvent(id string) (*models.Event, error, int) {
	query := `
		SELECT id, device_id, description, type, event, read FROM events
		WHERE id = $1
	`

	var event models.Event

	row := e.DB.QueryRow(query, id)
	err := row.Scan(
		&event.Id,
		&event.DeviceId,
		&event.Description,
		&event.Type,
		&event.Event,
		&event.Read,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, err, http.StatusNotFound
		default:
			return nil, err, http.StatusInternalServerError
		}
	}

	return &event, err, http.StatusOK
}

func (e EventModel) PutEvent(id string, data models.Event) (models.Event, error, int) {
	query := `
		UPDATE events SET read = $1
		WHERE id = $2
		RETURNING id, device_id AS deviceId, description, type, event, read
	`

	var event models.Event

	args := []interface{}{data.Read, id}
	row := e.DB.QueryRow(query, args...)
	err := row.Scan(
		&event.Id,
		&event.DeviceId,
		&event.Description,
		&event.Type,
		&event.Event,
		&event.Read,
	)

	if err != nil {
		return models.Event{}, err, http.StatusInternalServerError
	}

	return event, err, http.StatusOK
}

func (e EventModel) PostEvents(messagesChannel chan models.Event) (int, error) {
	stmt, err := e.DB.Prepare(`
			INSERT INTO events (device_id, description, type, event, read)
			VALUES ($1, $2, $3, $4, $5)
		`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	for event := range messagesChannel {
		if _, err := stmt.Exec(event.DeviceId, event.Description, event.Type, event.Event, event.Read); err != nil {
			log.Fatal(err)
		}
	}

	return len(messagesChannel), err
}

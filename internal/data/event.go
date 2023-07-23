package data

import (
	"database/sql"
	"errors"
	"net/http"
	"time"
)

type Event struct {
	Id          string    `json:"id"`
	DeviceName  string    `json:"deviceName"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Event       string    `json:"event"`
	Read        bool      `json:"read"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type EventModel struct {
	DB *sql.DB
}

func (e EventModel) GetEvents() ([]*Event, error) {
	query := `
	  SELECT id, device_name, description, type, event, read FROM events
	  ORDER BY updated_at
  	`

	rows, err := e.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []*Event

	for rows.Next() {
		var event Event

		err := rows.Scan(
			&event.Id,
			&event.DeviceName,
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

func (e EventModel) GetEvent(id string) (*Event, error, int) {
	query := `
		SELECT id, device_name, description, type, event, read FROM events
		WHERE id = $1
	`

	var event Event

	row := e.DB.QueryRow(query, id)
	err := row.Scan(
		&event.Id,
		&event.DeviceName,
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

func (e EventModel) PostEvent(data Event) (*Event, error, int) {
	query := `
		INSERT INTO events (device_name, description, type, event, read)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, device_name AS deviceName, description, type, event, read
	`

	var event Event

	args := []interface{}{data.DeviceName, data.Description, data.Type, data.Event, data.Read}
	// return the auto generated system values to Go object
	row := e.DB.QueryRow(query, args...)
	err := row.Scan(
		&event.Id,
		&event.DeviceName,
		&event.Description,
		&event.Type,
		&event.Event,
		&event.Read,
	)

	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return &event, err, http.StatusCreated
}

package data

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Event struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
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
	  SELECT * FROM events
	  ORDER BY updated_at`

	rows, err := e.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []*Event{}

	for rows.Next() {
		var event Event

		err := rows.Scan(
			&event.Id,
			&event.Name,
			&event.Description,
			&event.Type,
			&event.Event,
			&event.Read,
			&event.CreatedAt,
			&event.UpdatedAt,
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
		SELECT * FROM events
		WHERE id = $1`

	var event Event

	err := e.DB.QueryRow(query, id).Scan(
		&event.Id,
		&event.Name,
		&event.Description,
		&event.Type,
		&event.Event,
		&event.Read,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	invalidUuidErrorMessage := fmt.Sprintf(`pq: invalid input syntax for type uuid: "%s"`, id)

	if err != nil {
		switch {
		case err.Error() == invalidUuidErrorMessage:
			return nil, err, http.StatusBadRequest
		case errors.Is(err, sql.ErrNoRows):
			return nil, err, http.StatusNotFound
		default:
			return nil, err, http.StatusInternalServerError
		}
	}

	return &event, err, http.StatusOK
}

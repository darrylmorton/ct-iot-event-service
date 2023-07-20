package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Event struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Event       string `json:"event"`
	Read        bool   `json:"read"`
}

type EventsModel struct {
	Endpoint string
}

func (m *EventsModel) GetEvents() ([]Event, error) {
	resp, err := http.Get(m.Endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response []Event
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (m *EventsModel) GetEvent(id string) (*Event, error) {
	url := fmt.Sprintf("%s/%s", m.Endpoint, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response *Event
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

package test

import (
	"testing"
)

var db = DbConnection()
var eventNotFoundId = "00000000-0000-0000-0000-000000000000"

func getHealthCheckSuccess(t *testing.T) {
	expectedStatusCode := 200

	expectedResult := HealthCheck{Version: "0.0.1", Status: "available", Environment: "dev"}

	actualStatusCode, actualResult := GetHealthCheck()

	if expectedStatusCode != actualStatusCode {
		t.Errorf("expected: %v, actual: %v", expectedStatusCode, actualStatusCode)
	}

	AssertHealthCheck(expectedResult, actualResult, t)
}

func getEventsSuccess(t *testing.T) {
	expectedStatusCode := 200

	payload := CreateEventPayload()
	expectedResult, _ := CreateEvent(db, payload)

	actualStatusCode, actualResult := GetEvents()

	if expectedStatusCode != actualStatusCode {
		t.Errorf("expected: %v, actual: %v", expectedStatusCode, actualStatusCode)
	}

	AssertEvents([]Event{expectedResult}, actualResult, t)
}

func putEventSuccess(t *testing.T) {
	expectedStatusCode := 200

	_, results := GetEvents()

	expectedResult := results[0]
	expectedResult.Read = true

	actualStatusCode, actualResult := PutEvent(expectedResult.Id, expectedResult)

	if expectedStatusCode != actualStatusCode {
		t.Errorf("expected: %v, actual: %v", expectedStatusCode, actualStatusCode)
	}

	AssertEvent(expectedResult, actualResult, t)
}

func putEventInvalidUuid(t *testing.T) {
	expectedStatusCode := 400

	actualStatusCode, _ := PutEvent("1", Event{})

	if expectedStatusCode != actualStatusCode {
		t.Errorf("expected: %v, actual: %v", expectedStatusCode, actualStatusCode)
	}
}

func putEventNotFound(t *testing.T) {
	expectedStatusCode := 404

	actualStatusCode, _ := PutEvent(eventNotFoundId, Event{})

	if expectedStatusCode != actualStatusCode {
		t.Errorf("expected: %v, actual: %v", expectedStatusCode, actualStatusCode)
	}
}

func getEventSuccess(t *testing.T) {
	expectedStatusCode := 200

	_, expectedResult := GetEvents()

	actualStatusCode, actualResult := GetEvent(expectedResult[0].Id)

	if expectedStatusCode != actualStatusCode {
		t.Errorf("expected: %v, actual: %v", expectedStatusCode, actualStatusCode)
	}

	AssertEvent(expectedResult[0], actualResult, t)
}

func getEventInvalidUuid(t *testing.T) {
	expectedStatusCode := 400

	actualStatusCode, _ := GetEvent("1")

	if expectedStatusCode != actualStatusCode {
		t.Errorf("expected: %v, actual: %v", expectedStatusCode, actualStatusCode)
	}
}

func getEventNotFound(t *testing.T) {
	expectedStatusCode := 404

	actualStatusCode, _ := GetEvent(eventNotFoundId)

	if expectedStatusCode != actualStatusCode {
		t.Errorf("expected: %v, actual: %v", expectedStatusCode, actualStatusCode)
	}
}

func TestGroups(t *testing.T) {
	t.Run("Before", func(t *testing.T) {
		CreateDbTable(db)
		DeleteEvents(db)
	})

	t.Run("Health Check", func(t *testing.T) {
		t.Run("Success", getHealthCheckSuccess)
	})

	t.Run("Get Events", func(t *testing.T) {
		t.Run("Success", getEventsSuccess)
	})

	t.Run("Put Event", func(t *testing.T) {
		t.Run("Success", putEventSuccess)
		t.Run("Invalid UUID", putEventInvalidUuid)
		t.Run("Not Found", putEventNotFound)
	})

	t.Run("Get Event", func(t *testing.T) {
		t.Run("Success", getEventSuccess)
		t.Run("Invalid uuid", getEventInvalidUuid)
		t.Run("Not found", getEventNotFound)
	})

	t.Run("After", func(t *testing.T) {
		DropDbTable(db)
		db.Close()
	})
}

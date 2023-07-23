package test

import (
	"testing"
)

func postEventSuccess(t *testing.T) {
	expectedStatusCode := 201
	expectedResult := CreateEventPayload()

	actualStatusCode, actualResult := PostEvent()

	if expectedStatusCode != actualStatusCode {
		t.Errorf("expected: %v, actual: %v", expectedStatusCode, actualStatusCode)
	}

	AssertEvent(expectedResult, actualResult, t)
}

func getEventsSuccess(t *testing.T) {
	expectedStatusCode := 200
	expectedResult := []Event{CreateEventPayload()}

	actualStatusCode, actualResult := GetEvents()

	if expectedStatusCode != actualStatusCode {
		t.Errorf("expected: %v, actual: %v", expectedStatusCode, actualStatusCode)
	}

	AssertEvents(expectedResult, actualResult, t)
}

func getEventSuccess(t *testing.T) {
	expectedStatusCode := 200

	_, expectedResult := PostEvent()

	actualStatusCode, actualResult := GetEvent(expectedResult.Id)

	if expectedStatusCode != actualStatusCode {
		t.Errorf("expected: %v, actual: %v", expectedStatusCode, actualStatusCode)
	}

	AssertEvent(expectedResult, actualResult, t)
}

func getEventNotFound(t *testing.T) {
	expectedStatusCode := 404

	actualStatusCode, _ := GetEvent("00000000-0000-0000-0000-000000000000")

	if expectedStatusCode != actualStatusCode {
		t.Errorf("expected: %v, actual: %v", expectedStatusCode, actualStatusCode)
	}
}

func getEventInvalidUuid(t *testing.T) {
	expectedStatusCode := 400

	actualStatusCode, _ := GetEvent("1")

	if expectedStatusCode != actualStatusCode {
		t.Errorf("expected: %v, actual: %v", expectedStatusCode, actualStatusCode)
	}
}

func TestGroups(t *testing.T) {
	db := DbConnection()
	DeleteEvents(db)

	t.Run("Create Event", func(t *testing.T) {
		t.Run("Success", postEventSuccess)
	})

	t.Run("Get Events", func(t *testing.T) {
		t.Run("Success", getEventsSuccess)
	})

	t.Run("Get Event", func(t *testing.T) {
		t.Run("Success", getEventSuccess)
		t.Run("Invalid uuid", getEventInvalidUuid)
		t.Run("Not found", getEventNotFound)
	})
}
